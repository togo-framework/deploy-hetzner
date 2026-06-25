# deploy-hetzner — docs

**Hetzner deploy.** Provision a Hetzner Cloud server (cloud-init Docker) and run the app image.

## Install

```bash
togo install togo-framework/deploy-hetzner
```

Registers on the [`deploy`](https://github.com/togo-framework/deploy) base; select it with **deploy.provider in togo.yaml (or DEPLOY_PROVIDER)**, then use **`togo deploy`**.

## Interface

`Deployer` — `Provision`/`Deploy`/`Destroy`/`Status` over a `Spec{App,Dir,BuildCmd,Host,User,Image,Region,Domain}` built from your `togo.yaml`.

## Configuration

| Env var | Description |
|---|---|
| `HCLOUD_TOKEN` | Hetzner Cloud API token (required). |

## Usage & notes

Uses hcloud-go to create a server whose cloud-init runs `spec.Image`. `Destroy` deletes it.

## Example

```bash
togo deploy --provider hetzner --dry-run   # preview the plan
togo deploy --provider hetzner
```

## Links

- [hcloud-go](https://github.com/hetznercloud/hcloud-go)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/deploy-hetzner)
