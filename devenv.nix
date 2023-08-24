{ pkgs, ... }:

{
  # https://devenv.sh/packages/
  packages = [ pkgs.git pkgs.nodePackages.serverless ];

  # https://devenv.sh/languages/
  languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks = {
    nixpkgs-fmt.enable = true;
    revive.enable = true; # go lint https://github.com/mgechev/revive
    govet.enable = true;
    markdownlint.enable = true;
    nixfmt.enable = true;
    yamllint.enable = true;
  };

  pre-commit.settings = {
    markdownlint.config = {
      # MD013/line-length - Line length
      "MD013" = {
        # Number of characters
        "line_length" = 120;
      };
      yamllint.relaxed = true;
    };
  };
  enterShell = ''
    echo "[*]"
    FILE=.env
    if [ -f "$FILE" ]; then
      echo "[*] $FILE file exists. Remove it if you want to regenerate it."
    else
      echo "APP_ENV=dev" > .env
      echo "# Set GOROOT path for VC to load the correct Go version (from devenv)" >> .env
      echo "GOROOT=$GOROOT" >> .env
      echo "[*] $FILE generated."
    fi
    echo "[*]"
  '';

  # devenv creates .env and it doesn't use it
  dotenv.disableHint = true;

  # See full reference at https://devenv.sh/reference/options/
}
