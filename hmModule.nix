{pkgs, lib, config, ...}@args:
with lib;
let 
  nixgram = pkgs.callPackage ./package.nix {};
  dotenv = import ./.;
  nixgramService = dotenvFile: pkgs.writeShellScript "nixgram-service" ''
    ${dotenv}/bin/dotenv @${builtins.toString dotenvFile} -- ${nixgram}/bin/nixgram
  '';
in
{
  options = {
    services.nixgram = {
      enable = mkEnableOption "enable telegram interface for custom commands";
      dotenvFile = mkOption {
        type = with types; path;
        default = "/dev/null";
        example = "/some/file/with/the/credentials.env";
        description = "Dotenv file to load before starting nixgram";
      };
      customCommands = mkOption {
        type = with types; attrsOf str;
        default = {};
        example = {
          echo = "echo $*";
        };
        description = "Custom extra commands that the bot can handle";
      };
    };
  };
  config = mkIf config.services.nixgram.enable {
    systemd.user.services.nixgram = {
      Unit = {
        Description = "Command bot for telegram";
      };
      Service = {
        Type = "exec";
        ExecStart = "${nixgramService config.services.nixgram.dotenvFile}";
        Restart = "on-failure";
      };
      Install = {
        WantedBy = ["default.target"];
      };
    };
    home.packages = [
      nixgram
    ]
    ++ builtins.attrValues (
      builtins.mapAttrs 
        (name: value: pkgs.writeShellScriptBin ("nixgram-" + name) value)
        config.services.nixgram.customCommands
    );
  };
}
