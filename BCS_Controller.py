from ryu.base import app_manager
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.ofproto import ofproto_v1_3
from ryu.lib.packet import packet, ethernet
import json
import requests

BLOCKCHAIN_API = "http://blockchain-network/api/v1/reputation"

class BCSController(app_manager.RyuApp):
    OFP_VERSIONS = [ofproto_v1_3.OFP_VERSION]

    def __init__(self, *args, **kwargs):
        super(BCSController, self).__init__(*args, **kwargs)
        self.controller_id = "controller_1"  # Unique ID for this controller

    def get_reputation(self, controller_id):
        # Fetch reputation score from blockchain
        response = requests.get(f"{BLOCKCHAIN_API}/{controller_id}")
        if response.status_code == 200:
            return json.loads(response.text).get("reputation", 0)
        return 0

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def packet_in_handler(self, ev):
        msg = ev.msg
        datapath = msg.datapath
        ofproto = datapath.ofproto
        parser = datapath.ofproto_parser
        pkt = packet.Packet(msg.data)
        eth = pkt.get_protocol(ethernet.ethernet)

        # Reputation Check
        reputation_score = self.get_reputation(self.controller_id)
        if reputation_score < 0.5:
            self.logger.warning("Low reputation score, action required!")
            return

        # Handle normal packet forwarding
        in_port = msg.match['in_port']
        self.logger.info(f"Packet received on port {in_port}")
        actions = [parser.OFPActionOutput(ofproto.OFPP_FLOOD)]
        out = parser.OFPPacketOut(datapath=datapath, buffer_id=msg.buffer_id,
                                  in_port=in_port, actions=actions, data=msg.data)
        datapath.send_msg(out)
