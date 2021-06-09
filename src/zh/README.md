# TDengine 在 Kubenetes 上的部署

- 作者：Huo Linhe <lhhuo@taosdata.com>
- 更新日期：2021-06-09 16:24:00

为了支持[TDengine]在[Kubenetes][K8s]上的部署，特编写此文档。此文档完全开源，源码托管在 [taosdata/TDengine-Operator](https://github.com/taosdata/TDengine-Operator)，并欢迎所有人对此文档进行修改，您可以直接提交Pull Request，也可以添加 Issue，任何一种方式都将是我们的荣幸。TDengine完善离不开社区的共同努力，谢谢！

在本文档中，我们将从部署一套Kubenetes环境开始，介绍如何启动Kubenetes，并在Kubenetes上从头部署TDengine集群，简单介绍如何在K8s环境中进行TDengine集群的扩容和缩容，其中我们未能完整支持的地方也会有说明，可能出现问题的操作也作了简要的提示。

如果在实际操作过程中遇到问题，您总是可以通过官方微信 tdengine 联系到我们。

[TDengine]: https://github.com/taosdata/TDengine
[K8s]: https://kubernetes.io/
