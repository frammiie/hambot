import { ChatUserstate, Client } from 'tmi.js';
import { QuotesHandler } from './quotes.js';

export interface CommandHandler {
    regex: RegExp;
    handle: (client: Client, channel: string, tags: ChatUserstate, message: string, args: string[]) => void;
}

export const commands: CommandHandler[] = [
    QuotesHandler
];