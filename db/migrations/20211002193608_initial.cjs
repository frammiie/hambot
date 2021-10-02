exports.up = (knex) =>
    knex.schema
        .createTable('quotes', table => {
            table.increments('number').primary();
            table.string('content').notNull();
            table.string('author').notNull();
            table.string('by').notNull();
            table.datetime('added').notNull();
        });

exports.down = (knex) =>
    knex.schema
        .dropTable('quotes');