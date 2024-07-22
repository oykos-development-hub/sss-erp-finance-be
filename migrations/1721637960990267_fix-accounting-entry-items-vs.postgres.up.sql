alter table accounting_entry_items drop column date;

alter table accounting_entry_items add column date timestamp default NOW();

update accounting_entry_items set date = NOW();