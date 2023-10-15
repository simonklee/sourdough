#!/usr/bin/env python

import argparse
import sys
import sqlglot

def format_sql(query, dialect):
    return sqlglot.transpile(sql=query, 
                             write=dialect, # Target SQL dialect
                             max_text_width=100, # Maximum line width
                             normalize=True, # Normalize the query
                             pretty=True)   # Format the query

def main():
    parser = argparse.ArgumentParser(description='Format SQL queries using SQLGlot.')
    parser.add_argument('-f', '--file', help='File containing SQL queries. If not provided, reads from stdin.')
    parser.add_argument('-d', '--dialect', choices=['sqlite', 'postgres'], default='sqlite', help='Target SQL dialect for formatting.')
    parser.add_argument('-w', '--write', action='store_true', help='Update the file in place.')
    
    args = parser.parse_args()
    
    if args.write and not args.file:
        print("Error: -w option requires -f option to specify a file.", file=sys.stderr)
        sys.exit(1)
    
    if args.file:
        with open(args.file, 'r') as file:
            query = file.read()
    else:
        query = sys.stdin.read()
    
    formatted_queries = format_sql(query, args.dialect)
    formatted_output = ';\n\n'.join(formatted_queries) + ';'

    if args.write:
        with open(args.file, 'w') as file:
            file.write(formatted_output)
    else:
        print(formatted_output)

if __name__ == "__main__":
    main()
