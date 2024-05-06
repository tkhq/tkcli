class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.2"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "59e258836b3ac1d15efb49c2ff9637bcf8d72aa6a74165999525b0520af3c16c"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "da1534435bf06f6c988ec01c48c3a0645984327a349e26156c1d142df3d27ef4"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
