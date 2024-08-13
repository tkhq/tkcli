class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.3"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.3/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "01ab04b339665e62223015b7fcf8a8189956dc967967eaa7007856b69b1bf609"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.3/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "175db6f89e62eb192509600d63425333c43d283dbe2370bf2ca74a5b786fa808"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
