class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.3"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.3/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "f7689f57693b524f223b34175c0278267f74b6831464386038d0ee17ee546eff"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.3/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "07ffe737d5948b401e05a8d38fabedbd063f8ec69dda7d78084377ae65e5c014"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
