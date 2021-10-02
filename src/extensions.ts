import { ChatUserstate, Client } from 'tmi.js';

declare module 'tmi.js' {
    export interface ClientBase {
        respond: (channel: string, tags: ChatUserstate, message: string) => void;
        authorizedUser: (channel: string, tags: ChatUserstate) => boolean;
    }
}

Client.prototype.respond = function(this: Client, channel: string, tags: ChatUserstate, message: string) {
    this.say(channel, `@${tags.username} -> ${message}`);
}

const authorizedUsers = [ ...process.env.CHANNELS.split(' '), ...process.env.AUTHORIZED.split(' ')];
Client.prototype.authorizedUser = function(this: Client, channel: string, tags: ChatUserstate) {
    const authorized = tags.mod || authorizedUsers.includes(tags.username);
    if (!authorized) {
        this.respond(channel, tags, 'You are not allowed to run this command 🤡');
        return false;
    }
    return true;
}