#!/usr/bin/python3

import sys
import gtirb
from gtirb_capstone.instructions import GtirbInstructionDecoder


def extract_instruction_offsets(gtirb_file:str, insn_file:str):
    """
    Given a `gtirb_file`, extract all the instruction offsets
    and print them in a text file `insn_file`.
    """
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

extract_instruction_offsets(sys.argv[1], sys.argv[2])
