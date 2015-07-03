dog-cmp
==========


Program -> MainFunc StructDecl* FuncDecl*

MainFunc ->func main(){Statement}

StructDecl->type id struct{FieldDecl*}

FieldDecl -> id Type

VarDecl -> var id Type

FuncDecl -> func (id id) id (FormalList){VarDecl* Statement* return Exp}

FormalList
>　-> id Type FormalRest*

>　->

FormalRest -> , id Type

Type 
>　　-> int

>　　->bool

>　　->[]int

>　　->id

Statement 
> 　　-> {Statement*}

>　　->if Exp {Statement} else {Statement}

>　　->fmt.Println(Exp)

>　　->id := Exp

>　　->id = Exp

>　　->id[Exp] = Exp

>　　->for Exp {Statement*}


 Exp
> 　　-> Exp op Exp

>　　->Exp[Exp]

>　　->len(Exp)

>　　->Exp.id (ExpList)

>　　->INTERGER_LITERAL

>　　->true

>　　->false

>　　->id

>　　->[Exp]int

>　　->new(id)

>　　->id.id

>　　->!Exp

ExpList -> Exp ExpRest*

ExpRest -> ,Exp

###Lexer
38 kind of token
<table class="table table-bordered table-striped table-condensed">
   <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_AND</td>
      <td>&&</td>
   </tr>
 <tr>
      <td>TOKEN_ASSIGN</td>
      <td>=</td>
   </tr>
 <tr>
      <td>TOKEN_BOOL</td>
      <td>bool</td>
   </tr>
 <tr>
      <td>TOKEN_COMMER</td>
      <td>,</td>
   </tr>
 <tr>
      <td>TOKEN_DERIVE</td>
      <td>:=</td>
   </tr>
 <tr>
      <td>TOKEN_DOT</td>
      <td>.</td>
   </tr>
 <tr>
      <td>TOKEN_ELSE</td>
      <td>else</td>
   </tr>
 <tr>
      <td>TOKEN_FALSE</td>
      <td>false</td>
   </tr>
 <tr>
      <td>TOKEN_FMT</td>
      <td>fmt</td>
   </tr>
 <tr>
      <td>TOKEN_FOR</td>
      <td>for</td>
   </tr>
 <tr>
      <td>TOKEN_ID</td>
      <td>Identifier</td>
   </tr>
 <tr>
      <td>TOKEN_IF</td>
      <td>if</td>
   </tr>
 <tr>
      <td>TOKEN_INT</td>
      <td>int</td>
   </tr>
<tr>
      <td>TOKEN_LBRACE</td>
      <td>{</td>
   </tr>
<tr>
      <td>TOKEN_LBRACK</td>
      <td>[</td>
   </tr>
<tr>
      <td>TOKEN_LEN</td>
      <td>len</td>
   </tr>
<tr>
      <td>TOKEN_LPAREN</td>
      <td>(</td>
   </tr>
<tr>
      <td>TOKEN_LT</td>
      <td><</td>
   </tr>
<tr>
      <td>TOKEN_MAIN</td>
      <td>main</td>
   </tr>
<tr>
      <td>TOKEN_NEW</td>
      <td>new</td>
   </tr>
<tr>
      <td>TOKEN_NEWLINE</td>
      <td>\n</td>
   </tr>
<tr>
      <td>TOKEN_MAKE</td>
      <td>make</td>
   </tr>
<tr>
      <td>TOKEN_NOT</td>
      <td>!</td>
   </tr>
<tr>
      <td>TOKEN_NUM</td>
      <td>IntegerLiteral</td>
   </tr>
<tr>
      <td>TOKEN_PRINTLN</td>
      <td>Println</td>
   </tr>
<tr>
      <td>TOKEN_FUNC</td>
      <td>func</td>
   </tr>
<tr>
      <td>TOKEN_RBRACE</td>
      <td>}</td>
   </tr>
<tr>
      <td>TOKEN_RBRACK</td>
      <td>]</td>
   </tr>
<tr>
      <td>TOKEN_RETURN</td>
      <td>return</td>
   </tr>
<tr>
      <td>TOKEN_RPAREN</td>
      <td>)</td>
   </tr>
<tr>
      <td>TOKEN_SEMI</td>
      <td>;</td>
   </tr>
<tr>
      <td>TOKEN_START</td>
      <td>*</td>
   </tr>
<tr>
      <td>TOKEN_STRUCT</td>
      <td>struct</td>
   </tr>
<tr>
      <td>TOKEN_SUB</td>
      <td>-</td>
   </tr>
<tr>
      <td>TOKEN_TRUE</td>
      <td>true</td>
   </tr>
<tr>
      <td>TOKEN_TYPE</td>
      <td>type</td>
   </tr>
<tr>
      <td>TOKEN_VOID</td>
      <td>void</td>
   </tr>
</table>




