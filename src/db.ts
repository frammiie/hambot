import knex from 'knex';
import config from '../db/config.json';

export interface Quote {
    number: number;
    content: string;
    author: string;
    added: Date;
    by: string;
}

export const db = knex(config);