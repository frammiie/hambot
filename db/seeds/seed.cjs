exports.seed = async (knex) => {
    await knex('quotes').del();
    await knex('quotes').insert([
        {
            number: 1, content: 'I\'m quote #1!', author: 'frammie', by: 'seed', added: new Date()
        },
        {
            number: 2, content: 'I\'m another quote with number #2!', author: 'frammie (not really)', by: 'seed', added: new Date()
        }
    ])
}