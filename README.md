A simple generative fuzzer
--------------------------

Sample grammar:

```
START := paren ;
paren := "(" opt_paren_list ")" ;
opt_paren_list := paren_list | EMPTY ;
paren_list := paren_list paren | paren ;
```
