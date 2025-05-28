class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.4"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.4/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "e91b356d6e79d5788e4db9a2bedf7bee6432933b2a640b931596ea1ee79bce30"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.4/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "e5ef1e881879d1f53514f82860b3cf6d0b2f3d8decb3ef68831606552156e516"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
