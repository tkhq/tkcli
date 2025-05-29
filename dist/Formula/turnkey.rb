class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.5"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.5/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "729804cf6652e23b8e3fc0a2548e0e0327b69826f50b10b560656f77af76f575"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.5/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "36ddc3f9675214c35e924f8212028e35d7674ae1e6a46a49d68fa4b983c1d954"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
