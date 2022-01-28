# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Hammer < Formula
  desc ""
  homepage ""
  version "0.12.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal/hammer/releases/download/v0.12.0/hammer_darwin_amd64.tar.gz"
      sha256 "2211d9639819a34e18f5d50947ac5f8a3ad192e24695f32d52a9c65a8bed1328"

      def install
        bin.install "hammer"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal/hammer/releases/download/v0.12.0/hammer_linux_amd64.tar.gz"
      sha256 "f7e066017e5f59dd5c8ee5c2cd4f2e4e8333e3a45bdd574cc26a1cc4e950aa81"

      def install
        bin.install "hammer"
      end
    end
  end

  test do
    system "#{bin}/hammer version"
  end
end
