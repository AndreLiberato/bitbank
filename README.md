# bitbank

## Integrantes 

> André Luiz de Sena Liberato

> Rian Abdias Balbino de Azevedo

## Tecnologias Utilizadas

> Go Lang 

## CI/CD

Este projeto utiliza GitHub Actions para executar a rotina de integração contínua.
Seguindo o fluxo adotado no projeto, a branch `main` exerce o papel de branch de integração.

A pipeline executa automaticamente as seguintes etapas:

> Resolução das dependências

> Build

> Execução dos testes unitários

> Criação de uma tag no formato `build-XXX` para identificar a configuração gerada

Além da pipeline principal, o projeto possui uma pipeline de estabilização para branches `stabilization/rc-*`.
Ao aceitar pull requests de `main` ou `bugfix*` nessas branches, a automação executa dependências, build, análise estática com `go vet`, testes unitários, cria a tag `rc-*` correspondente e disponibiliza um pacote `.zip` da versão na área de artifacts do workflow.
