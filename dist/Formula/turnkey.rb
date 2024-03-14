class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.0"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.0/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "b3c53d7a2ce8a99cd5cfc39f4f17fd7520c85ede6d311e14c850878d341bb666"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.0/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "b3c53d7a2ce8a99cd5cfc39f4f17fd7520c85ede6d311e14c850878d341bb666"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
