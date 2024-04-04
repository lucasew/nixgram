{ buildGoModule
, lib
}:

buildGoModule {
  name = "nixgram";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-cmiQg6EEbaJDTzTt+uHMo8Y3lH4MiZEnpn1Ef8hbbpA=";

  meta = with lib; {
    description = "A smol software that can interface with commands using telegram";
    homepage = "https://github.com/lucasew/nixgram";
    platforms = platforms.linux;
  };
}
