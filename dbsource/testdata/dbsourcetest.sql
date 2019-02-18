# This file creates an in-memory database to be used for testing.

create table company(id string, name string);

create table person(id string, companyid string, firstname string, lastname string);

insert into company(id, name) values
  ('x', 'Example, Inc.'),
  ('s', 'Sample Corp.');

insert into person(id, companyid, firstname, lastname) values
  ('p1', 'x', 'John', 'Doe'),
  ('p2', 'x', 'Jane', 'Doe'),
  ('p3', 's', 'Tom', 'Smith'),
  ('p4', 's', 'Tim', 'Smith'),
  ('p5', 's', 'Jim', 'Smith');
