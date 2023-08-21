{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = [
    pkgs.git
    pkgs.nodePackages.serverless
  ];

  # https://devenv.sh/languages/
  languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  pre-commit = {
  };

  # https://devenv.sh/scripts/
  # scripts.deploy.exec = "sls deploy";

  enterShell = ''
    echo "---"
    git --version
    go version
    echo "Serverless Framework"
    sls --version
    echo "---"
  '';

  # See full reference at https://devenv.sh/reference/options/
}
