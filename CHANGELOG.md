# Changelog

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
