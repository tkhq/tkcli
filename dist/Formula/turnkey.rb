class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.4"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "356b4e1636a2f6cea0242907b451098a6d3c939d272901a4801458d209f931a9"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "68ffb9e5dbb0945b5fd53a2e2ccadfdba1f620720865c23e275fc159f70d30de"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
