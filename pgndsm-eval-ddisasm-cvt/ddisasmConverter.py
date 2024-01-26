#!/usr/bin/python3

import sys
import gtirb
from gtirb_capstone.instructions import GtirbInstructionDecoder


def convert(gtirb_file, insn_file):
    ir = gtirb.IR.load_protobuf(gtirb_file)
    module = ir.modules[0]
    decoder = GtirbInstructionDecoder(module.isa)
    insns = []
    for block in module.code_blocks:
        insns.extend((insn.address for insn in decoder.get_instructions(block)))
    insns.sort()
    with open(insn_file, "w") as f:
        for addr in insns:
            print(hex(addr), file=f)


if len(sys.argv) < 3:
    print(f"expected two arguments, given {len(sys.argv)}")
    exit(1)

convert(sys.argv[1], sys.argv[2])
