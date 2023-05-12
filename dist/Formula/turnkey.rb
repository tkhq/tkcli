class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.1"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.1/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "8cbb2d128ebaa86fe579faf5953a82798dfb72cc697879ef0db2738fec28e295"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.1/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "5f8ec155a69c780c0cbb682d18c43fd74d1a47964a64e18d4bf10c7057cb53e5"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
