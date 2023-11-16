class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.5"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.5/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "03ed1ce9fbe2ae91f3f78fea1051459130b917746c854fb7856604a605587acf"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.5/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "a9721fc5dc7edd5de60740a60113d8a7195392db5f81c0240898ab2c9cb46fd7"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
