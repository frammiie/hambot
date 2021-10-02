import { ChatUserstate, Client } from 'tmi.js';
import { CommandHandler } from '.'
import { db, Quote } from '../db.js';

const escapeRegex = new RegExp('\\\\', 'g');

export const QuotesHandler: CommandHandler = {
    regex: new RegExp('quotes?'),
    handle: async (client, channel, tags, message, args) => {
        // TODO: Multilevel regex options
        // TODO: Permission model
        switch(args[0]) {
            case 'add':
                if (!client.authorizedUser(channel, tags)) return;

                if (args.length !== 3) {
                    client.respond(channel, tags, 'Format: "[content]" [author] ℹ️');
                    return;
                }
                const addedNumber = await db<Quote>('quotes').insert({
                    content: args[1].substring(1, args[1].length-1).replaceAll(escapeRegex, ''),
                    author: args[2],
                    added: new Date(),
                    by: tags.username,
                }).returning('number');

                client.respond(channel, tags, `Added quote #${addedNumber} successfully 📝`);
                break;
            case 'del':
                if (!client.authorizedUser(channel, tags)) return;

                const delNum = parseNumber(client, channel, tags, args[1]);
                if (!delNum) return;
                await db<Quote>('quotes').del().where('number', delNum);
                client.respond(channel, tags, `Deleted quote #${delNum} successfully 💀`);
                break;
            default:
                // Assume number request
                const number = parseNumber(client, channel, tags, args[0]);
                if (!number) return;
                const quote = await db<Quote>('quotes').where('number', number).first();
                if (!quote) {
                    client.respond(channel, tags, `Quote with number ${number} not found! 👀`);
                    return;
                }
                client.say(channel, `"${quote.content}" - ${quote.author}`);
        }
    }
};

const parseNumber = (client: Client, channel: string, tags: ChatUserstate, input: string) => {
    const number = parseInt(input);
    if (isNaN(number)) {
        client.respond(channel, tags, `Please enter a number of a quote. 🔢`);
        return false;
    }
    return number;
};