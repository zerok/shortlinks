package main

import (
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging/mem"
)

var _ = pkger.Apply(mem.UnmarshalEmbed([]byte(`1f8b08000000000000ffec97ff73a2b81bc7ff959de7d70fab090808339f1faaad96edea5614ac763a9d0829727c0947705bdbf17fbf0928ea5eeb766f76e6eee6761c8624bcdf499e27c90b7c81307d601ccc1710d7799883094d87d39c379f69cea226cfbd661016cbd5a2e1b164d7b86479118769c441022bc9585e5c93620926800443925030e14873cebceae184e4012daab2cdd8b6342085b704335dc5b104e382c414cc071273baadd994709656da3eeb8531e53b75357a5d3da7595d9e505e7ca3164ddf3806cc5f89f15e601bc2e96807244cc12cf21595fe5abafa6cc0fc1fb43503d648985fba5d9af3b04c066e6005369b8d040f558c2fa7e66e369330c84911b2b4948a8517779f16248ccba6b45ab9039d043c7ca6601a9a0409f32998326ee9ad760bb770d9725f84a54746d8f888e58fb23cc1c8543413a18666c89aa2ab6de57f089b08810421bff745ceaaf4f17539ea39fd0a26d6745d976554ae0715f5b6ae4b301453075391c04a1998aa6e20cd109371421f4c156109faa2248cb62f3a42125c13ffde0bd83d02f31649e5ef4e823331539e514f8c391637acea3a42baa219120cb96869b5b18e75241b1b09066febd14e5f47b891a0fb83fa4e9817cbd73c9ad1920d55db7914d9406dac607923c1b85e8c4eccbc8897e176e2a85aa416124f7a3109aa077d9a96f7cf3c23392d8ba36d51e4e5aedc3a4b9aff9c937f5aba8542224ecf6b38381a7911a6c7ee2d290e457b6adc9e1cfaae064ab5e98e79c245ed834f339afa34f5d6e6877c95965bfa0034b7904501f51b0113bdbd0a9c5bd8f9eefe33dc91202b27fe02d751f0de0d7008a18d043e29c82e216267a6c5bebfbda91cec9d686b2284f07d98864548e2c62a6bf0dfe3d3b87bcdb0e39e2cb776e06b09c67c9778584758d3f43f116ffb363b853c45c1ed768d3cbc479eacab2afe89c8c3a8add648c208894864e5ddc8db86f86ee4edf56f224f7810aa3daa6a60dd30b05623af5c871df3da7f1bf39a47afc6373e7c8e34bf3e7cfeb9007a1515359460843ff546aeddb37a7667143df5acfe53b6488a67ab1b5c59ddb3a0bcfa71645dce975e3264b31b8b0d7e3bcbaccbce9a4cf1d24ba380c8ae6af5d5af7eb7b3f2a74ffcf3190b765effe613b72e87c84be2d57cdd592d9451b0485dbee83e1ef43f5ccfa73d349fda0fb39b51e0f7e3627e3344b3297e3cecd7eadbf13ce9e1c5e5289829eeda4bdc95df331099e2d853ece5e2e6ecea6ac2afbad1706c3b3dc71e77dc4914dbaee3046347bdb0dd20183bb36012198e75e1cec6eed071d61d4ea66aee61037bc9d30399da2debc2f8b2cfc3127bc963f6e591fd1fbe83efead0d15fff31fea5476df307000000ffff010000ffff688afc36290e0000`)))
