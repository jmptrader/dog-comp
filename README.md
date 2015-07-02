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
<table class="table table-bordered table-striped table-condensed">
   <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
 <tr>
      <td>TOKEN_ADD</td>
      <td>+</td>
   </tr>
</table>




