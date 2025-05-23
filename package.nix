{ buildGoModule
, lib
}:

buildGoModule {
  name = "nixgram";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-QoVryPstnfdvc390pgAeLF04AXNzR2BEXLXY2isTK+s=";

  meta = with lib; {
    description = "A smol software that can interface with commands using telegram";
    homepage = "https://github.com/lucasew/nixgram";
    platforms = platforms.linux;
  };
}
