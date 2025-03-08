# A default development shell for the Python API server.
# If outputs are generated in a separate shell (direnv),
# activate the Python virtual environment manually.

{ pkgs ? import <nixpkgs> {}}:

pkgs.mkShell{
  packages = with pkgs; [
    (python3.withPackages (python-pkgs: [
      python-pkgs.pip
    ]))
  ];

  shellHook = ''
    python -m venv .venv
    source .venv/bin/activate
    pip install -r ./requirements.txt
  '';
}
