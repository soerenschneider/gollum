# Changelog

## [1.2.1](https://github.com/soerenschneider/gollum/compare/v1.2.0...v1.2.1) (2025-07-07)


### Bug Fixes

* **deps:** Update dependency go to v1.24.4 ([274b5b4](https://github.com/soerenschneider/gollum/commit/274b5b46d20f821bd907a9a57738c48d488c67e3))
* **deps:** Update dependency go to v1.24.4 ([7c8a65b](https://github.com/soerenschneider/gollum/commit/7c8a65b4ac34b20e8fc4e2ae90a6ad540ab65671))
* **deps:** Update golang.org/x/exp digest to b7579e2 ([d8037d6](https://github.com/soerenschneider/gollum/commit/d8037d659204dd3b4ef5bb16dd245a7376500ebf))
* **deps:** Update golang.org/x/exp digest to b7579e2 ([3704019](https://github.com/soerenschneider/gollum/commit/370401984b924081a09470f3299f07159375e787))
* **deps:** Update k8s.io/utils digest to 4c0f3b2 ([dfe50d2](https://github.com/soerenschneider/gollum/commit/dfe50d2dbfdf932355d91224cc1c4fd7a95cda88))
* **deps:** Update k8s.io/utils digest to 4c0f3b2 ([01de913](https://github.com/soerenschneider/gollum/commit/01de913284d6a699a883ea1c2f7332b098044085))
* **deps:** Update knative.dev/pkg digest to 16de760 ([b731790](https://github.com/soerenschneider/gollum/commit/b731790f59aa8620285efd428575c0122ba4a84e))
* **deps:** Update knative.dev/pkg digest to 16de760 ([5ee0512](https://github.com/soerenschneider/gollum/commit/5ee05123c15f1dc7b5f0cf99d32e51c93e20cdcd))
* **deps:** Update module github.com/google/cel-go to v0.25.0 ([885dae7](https://github.com/soerenschneider/gollum/commit/885dae719a52bd35d89c3dc6c27e3bdb4ec1229b))
* **deps:** Update module github.com/google/cel-go to v0.25.0 ([c33fbbd](https://github.com/soerenschneider/gollum/commit/c33fbbdeaa4d154c110a4218621fdc19f55236ec))
* **deps:** Update module github.com/hashicorp/go-retryablehttp to v0.7.8 ([764eb59](https://github.com/soerenschneider/gollum/commit/764eb5960b197482353a8e10863ed2db947736ec))
* **deps:** Update module github.com/hashicorp/go-retryablehttp to v0.7.8 ([1a5bb5b](https://github.com/soerenschneider/gollum/commit/1a5bb5bf5a714aa4af41629e555a977b1c0c2d96))
* **deps:** Update module github.com/Masterminds/semver/v3 to v3.4.0 ([ce77205](https://github.com/soerenschneider/gollum/commit/ce7720564af9f1081f61f8f99c60e4af096229a8))
* **deps:** Update module github.com/Masterminds/semver/v3 to v3.4.0 ([b1cd9de](https://github.com/soerenschneider/gollum/commit/b1cd9deeaff4afc84c7b6c0bce5c699145aaf8c1))
* **deps:** Update module github.com/tektoncd/pipeline to v1 ([9fa6e2c](https://github.com/soerenschneider/gollum/commit/9fa6e2cee0f0fe6979b0b889b5bc6d1a9257e2c5))
* **deps:** Update module github.com/tektoncd/pipeline to v1 ([4bb9008](https://github.com/soerenschneider/gollum/commit/4bb9008d8492fb39b3e560969224334b71c22d8e))
* **deps:** Update module k8s.io/api to v0.33.2 ([e3288fc](https://github.com/soerenschneider/gollum/commit/e3288fc4e65d2d208f0a36ccea6e64d9b9fe4cdb))
* **deps:** Update module k8s.io/api to v0.33.2 ([6602de8](https://github.com/soerenschneider/gollum/commit/6602de89b45778b41879072445017befa028c2f0))
* **deps:** Update module k8s.io/apimachinery to v0.33.2 ([2c449d3](https://github.com/soerenschneider/gollum/commit/2c449d30bdf7b6f71613ef55d292fcdf5a710366))
* **deps:** Update module k8s.io/apimachinery to v0.33.2 ([c4f234c](https://github.com/soerenschneider/gollum/commit/c4f234c44edb6275932a54990ffc91a7914ac475))
* **deps:** Update module k8s.io/client-go to v0.33.2 ([d5abdee](https://github.com/soerenschneider/gollum/commit/d5abdeea62ed955bc7934107f44d09aaa01d3f57))
* **deps:** Update module k8s.io/client-go to v0.33.2 ([dc47a30](https://github.com/soerenschneider/gollum/commit/dc47a30ca6e72579167b68a86430be18e5c246b7))
* **deps:** Update module sigs.k8s.io/controller-runtime to v0.21.0 ([b2afd2f](https://github.com/soerenschneider/gollum/commit/b2afd2f1cfb44cdfc14b5e7159cf6a49f0bde8ac))
* **deps:** Update module sigs.k8s.io/controller-runtime to v0.21.0 ([2826500](https://github.com/soerenschneider/gollum/commit/2826500431cf8f0a85f1f52856f1368934a01401))
* fix panic if value for repo or owner is too short (&lt;5) ([00f10cb](https://github.com/soerenschneider/gollum/commit/00f10cbfc5d303a7b36c32c591f38e4014a83ebf))

## [1.2.0](https://github.com/soerenschneider/gollum/compare/v1.1.0...v1.2.0) (2025-04-16)


### Features

* add flags to control requeueing and logging ([66d651b](https://github.com/soerenschneider/gollum/commit/66d651bc61a02e7ee65a2835a24c28d26d9103a9))


### Bug Fixes

* **deps:** bump github.com/prometheus/client_golang ([4cd1939](https://github.com/soerenschneider/gollum/commit/4cd1939aafc1dc925f66ee20d4a5958ad5021e83))
* **deps:** bump github.com/prometheus/client_golang from 1.20.5 to 1.22.0 ([e3b7c36](https://github.com/soerenschneider/gollum/commit/e3b7c367b7340f85becf4e7782f85be8bf7b7317))
* **deps:** bump github.com/tektoncd/pipeline from 0.68.0 to 0.70.0 ([4658720](https://github.com/soerenschneider/gollum/commit/4658720301b2ea8fe8c50371ee127fceb57f43f0))
* **deps:** bump github.com/tektoncd/pipeline from 0.68.0 to 0.70.0 ([6ed6232](https://github.com/soerenschneider/gollum/commit/6ed6232dc6817a97bc73c4552040bf0b5081d895))
* **deps:** bump golang from 1.24.0 to 1.24.2 ([3eef829](https://github.com/soerenschneider/gollum/commit/3eef8294e2ca35b933cd0d7cb02a9cee73f67124))
* **deps:** bump golang from 1.24.0 to 1.24.2 ([dac8423](https://github.com/soerenschneider/gollum/commit/dac84239c33557af49e7343f6bf818bf7633ae08))
* **deps:** bump golang.org/x/net from 0.34.0 to 0.36.0 ([72ad83c](https://github.com/soerenschneider/gollum/commit/72ad83c86959c19eef8bf4188d965f839d019179))
* **deps:** bump golang.org/x/net from 0.34.0 to 0.36.0 ([c12785c](https://github.com/soerenschneider/gollum/commit/c12785c5258b29fcef9b3b48f2ed3a37d541a183))
* **deps:** bump k8s.io/api from 0.32.2 to 0.32.3 ([b0dd8f3](https://github.com/soerenschneider/gollum/commit/b0dd8f3767a6d8adba8e9a23e246da5e6909da9b))
* **deps:** bump k8s.io/api from 0.32.2 to 0.32.3 ([280d0a6](https://github.com/soerenschneider/gollum/commit/280d0a6e345c050df1a4011e49eac65344e19dba))
* **deps:** bump sigs.k8s.io/controller-runtime from 0.19.1 to 0.20.4 ([e090d17](https://github.com/soerenschneider/gollum/commit/e090d1734b1f4d98f96d840bdc100615914846f4))
* **deps:** bump sigs.k8s.io/controller-runtime from 0.19.1 to 0.20.4 ([297c12f](https://github.com/soerenschneider/gollum/commit/297c12feab2d1eafd1e6a2ddead305f6a1c19a7e))
* fix logic to only process releases during 'on' hours ([11dd84e](https://github.com/soerenschneider/gollum/commit/11dd84e3f027924f1236e70cfc0b7c164734eee9))
* inverse result after recent refactoring ([c996224](https://github.com/soerenschneider/gollum/commit/c996224b781dd1fd49228ad681424367caed53ce))
* register metrics ([32e605f](https://github.com/soerenschneider/gollum/commit/32e605fa05aea1ab7870e7ae7a4215b862b29741))

## [1.1.0](https://github.com/soerenschneider/gollum/compare/v1.0.0...v1.1.0) (2025-02-18)


### Features

* feature to omit certain versions ([072cbfb](https://github.com/soerenschneider/gollum/commit/072cbfb70b097f79a0cb4ff56410dcadec208577))


### Bug Fixes

* **deps:** bump github.com/tektoncd/pipeline from 0.65.2 to 0.68.0 ([866f948](https://github.com/soerenschneider/gollum/commit/866f9487d0dc117b86037540e2ea0a4499baafe0))
* **deps:** bump golang from 1.23.4 to 1.24.0 ([ae7794d](https://github.com/soerenschneider/gollum/commit/ae7794dba6e6a004f55568f0a71e444365e3e122))
* **deps:** bump k8s.io/api from 0.32.0 to 0.32.2 ([d112294](https://github.com/soerenschneider/gollum/commit/d112294707356cf21326a9e8ea89e0bf94bb4523))
* **deps:** bump k8s.io/apimachinery from 0.32.0 to 0.32.2 ([d66e277](https://github.com/soerenschneider/gollum/commit/d66e277a64c86ab6a7ddb335162eda6da86e6cf5))
* **deps:** bump k8s.io/client-go from 0.32.0 to 0.32.2 ([35a1144](https://github.com/soerenschneider/gollum/commit/35a11445b720f5e22d38dc61ad5ac3fa82cb6174))

## 1.0.0 (2025-01-09)


### Bug Fixes

* **deps:** bump github.com/prometheus/client_golang ([e8d5d51](https://github.com/soerenschneider/gollum/commit/e8d5d51f972cc8a832151c79dd7e79d5e4f62fcc))
* **deps:** bump golang.org/x/net from 0.28.0 to 0.33.0 ([6c554b5](https://github.com/soerenschneider/gollum/commit/6c554b5474e97bcfaabd0b7f9e925bcff873caca))
* **deps:** bump google.golang.org/api from 0.181.0 to 0.216.0 ([6de70b4](https://github.com/soerenschneider/gollum/commit/6de70b4a932aedf542cff89d3b485f7c4b4eb554))
* **deps:** bump k8s.io/client-go from 0.31.0 to 0.32.0 ([e1d0b04](https://github.com/soerenschneider/gollum/commit/e1d0b043ce877c70bcec0f3f93e9891128a1f4d3))


### Miscellaneous Chores

* release 1.0.0 ([d69e3ac](https://github.com/soerenschneider/gollum/commit/d69e3ac54863803782b329e88519efacb9fae091))
