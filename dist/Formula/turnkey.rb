class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.2"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.2/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "c5cba5ecebe7ac394afaffc2e171708369e829692c3e2efaeb60ada236c8d974"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.2/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "8d5f1d5315db97801f9266fa60fb45ca600b2e95b1782fdd7886c23bfb109fed"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
