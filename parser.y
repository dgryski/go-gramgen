%{
package main

import (
    "log"
)

var symtab = newSymbolTable()

%}

%union {
    s string
    g generator

    seqs *sequence
    choices *choice
}

%token <s> tID tQSTRING

%token tASSIGN tDOTDOT tQSTRING

%type <s> rule rule_list

%type <g> expr
%type <seqs> expr_seq
%type <choices> expr_list

%%

grammar : rule_list ;

rule_list : rule_list rule 
    | rule
    ;

rule : tID tASSIGN expr_list ';' {
    if err := symtab.addRule($1, $3); err != nil {
        log.Fatalln("error: ", err)
    }

    $$ = $1
}

expr_list : expr_list '|' expr_seq { $1.add($3); $$ = $1; }
    | expr_seq { $$ = &choice{c:[]generator{$1}} }
    ;

expr_seq : expr_seq expr { $1.add($2); $$ = $1; }
    | expr { $$ = &sequence{s:[]generator{$1}} }
    ;

expr: tQSTRING {  $$ = terminal($1) }
    | tID { $$ = symtab.addVariable($1) }
    ;

