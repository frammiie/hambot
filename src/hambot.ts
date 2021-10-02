import { Client } from 'tmi.js';
import { commands } from './commands/index.js';
import './extensions.js';

const client = new Client({
    connection: {
        reconnect: true,
        secure: true,
    },
    identity: {
        username: 'frammiebot',
        password: process.env.TOKEN
    },
    channels: process.env.CHANNELS.split(' '),
});

console.info('frammiebot is starting...');

client.connect()
    .catch(console.error)
    .then(() => console.info('frammiebot connected successfully'));

client.on('join', (channel) => client.say(channel, 'frammiebot is here!'));

const argsRegex = new RegExp('"[^"]*"|[^ ]+', 'g');
client.on('message', async (channel, tags, message, self) => {
    if (self) return;
    if (!message.startsWith('!')) return;
    let args = message.substring(1).match(argsRegex);
    const command = args[0];
    args = args.slice(1);
    for(let handler of commands) {
        if (handler.regex.test(command))
            handler.handle(client, channel, tags, message, args);
    }
});