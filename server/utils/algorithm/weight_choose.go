package algorithm

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	CurIndex int
	Rss      []*WeightNode
	Rsw      []int
}
type WeightNode struct {
	Addr            string //服务器地址
	Weight          int    //权重值
	CurrentWeight   int    //节点当前权重
	EffectiveWeight int    //有效权重
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	//这里拿到权重
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	//实例化具体的Node节点
	node := &WeightNode{Addr: params[0], Weight: int(parInt)}
	node.EffectiveWeight = node.Weight //权重值=有效权重
	r.Rss = append(r.Rss, node)        //append到服务器节点
	return nil
}

//获取
func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode //该次最优的ip
	for i := 0; i < len(r.Rss); i++ {
		w := r.Rss[i]
		//统计所有有效权重之和
		total += w.EffectiveWeight
		//变更节点临时权重为的节点临时权重+节点有效权重
		w.CurrentWeight += w.EffectiveWeight
		//有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}
		//选择最大临时权重点节点
		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	//变更临时权重为 临时权重-有效权重之和
	best.CurrentWeight -= total
	return best.Addr
}
