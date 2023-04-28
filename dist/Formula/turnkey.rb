class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.0"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.0/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "5e2367343db0b99d6074db64d0a0f1e3e09e982f5e9a4e375b4fcdf8a3424359"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.0/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "08d7dbbec8fcb860aa796bd5c281ce97e04d7348fdce215704faad0eadf240a7"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
