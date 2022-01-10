#!/usr/bin/env bash

set -e

# Reinitialize the database.
psql -d test -c 'drop schema public cascade'
psql -d test -c 'create schema public'
psql -d test -c 'create extension "uuid-ossp"'
psql -d test -c 'alter schema public owner to de'
