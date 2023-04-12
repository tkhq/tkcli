class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v0.3.3"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v0.3.3/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "dab249ca4b1bd70ac92b1e403c68b44b4aafc96b76524be1a4484013566ac182"

    def install
      bin.install "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v0.3.3/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "52ff1db2e8276f5a3adebf72b790dc3907490f317c1886e0b88af30fb958200b"

    def install
      bin.install "turnkey"
    end
  end

end
