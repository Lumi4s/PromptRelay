{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  packages = with pkgs; [
    go
    gopls          # языковой сервер для автокомплита в редакторе
    golangci-lint  # линтер
    delve          # отладчик
  ];
}
