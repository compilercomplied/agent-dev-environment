# Python SDK Action

This action handles the generation, building, and optional publishing of a Python SDK from an OpenAPI specification.

## Placeholders used in this guide
- `<PACKAGE_NAME>`: The name of the package (e.g., `myorg-my-repo`).
- `<PACKAGE_MODULE>`: The importable name (usually `<PACKAGE_NAME>` with hyphens replaced by underscores).
- `<VERSION>`: The generated version string, including the beta suffix (e.g., `1.0.42b7`).

---

## Consuming Beta Packages (Pull Requests)

To ensure SDK changes are safe before merging, the CI pipeline generates **Beta SDKs** for every Pull Request. These are uploaded as GitHub Artifacts rather than being published to the live PyPI registry.

### 1. Download the Artifact
1. Navigate to your **Pull Request** on GitHub.
2. Click the **Actions** tab (or the "Details" link on the `openapi-sdks` job).
3. Scroll to the **Artifacts** section at the bottom of the run summary.
4. Download the artifact named `sdk-python-<VERSION>`.

### 2. Install Locally
Unzip the downloaded artifact and use `pip` to install the package from the `dist/` directory:

```bash
# Install the Wheel file (recommended)
pip install ./path/to/extracted/dist/<PACKAGE_NAME>-<VERSION>-py3-none-any.whl

# Or install the Source distribution
pip install ./path/to/extracted/dist/<PACKAGE_NAME>-<VERSION>.tar.gz
```

### 3. Verify the Installation
Confirm the beta version is active in your environment:

```bash
python -c "import <PACKAGE_MODULE>; print(<PACKAGE_MODULE>.__version__)"
```

---

## Internal Parameters
To maintain a clean parent workflow, the output and distribution paths are defined as **internal parameters** with default values at the top of the `action.yaml` file. Change them there if you need to restructure the internal generation path.
