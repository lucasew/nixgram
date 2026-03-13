{ buildGoModule
, lib
}:

buildGoModule {
  name = "nixgram";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256:0000000000000000000000000000000000000000000000000000000000000000";

  meta = with lib; {
    description = "A smol software that can interface with commands using telegram";
    homepage = "https://github.com/lucasew/nixgram";
    platforms = platforms.linux;
  };
}
