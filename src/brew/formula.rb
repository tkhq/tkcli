class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "$VERSION"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/$VERSION/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "$DARWIN_AARCH64_SHA256"

    def install
      bin.install "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/$VERSION/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "$DARWIN_X86_64_SHA256"

    def install
      bin.install "turnkey"
    end
  end

end
