let
  pkgs = import <nixpkgs> {};
in pkgs.buildGoModule rec {
  name = "nixgram";
  version = "0.0.1";
  vendorSha256 = "143fbg47yi3xlqkr328cgsa3gim3rkhzmv9l9x1s4v84l61r0s3j";
  src = ./.;
  meta = with pkgs.lib; {
    description = "A smol software that can interface with commands using telegram";
    homepage = "https://github.com/lucasew/nixgram";
    platforms = platforms.linux;
  };
}
