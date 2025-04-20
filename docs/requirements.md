# D&D 5e Converter

As a D&D Dungeon Master you need the ability to quickly search and find
information during a game. We achieve this by accessing a large database of
Markdown files containing detailed information about spells, items and
creatures.

All the necessary information is already stored locally in multiple JSON files.
These must be parsed and converted to Markdown. This is done using a Go program
that reads the JSON files from another local repository and then writes the
Markdown files to a local directory.

## Data Location and Format

For each of the main categories (spells, creatures and items) there is a folder
and an associated index file, e.g. `<repo>/data/spells/index.json` file. This
file contains the locations to all other relevant JSON files that are also
located in the same folder. All the referenced JSON files contain a single
object with a single child that is the actual data in array format.

Below are some examples of the JSON files associated with spells. All other
files follow a similar format, but not necessarily the same.

### index.json file for spells

```json
{
    "AAG": "spells-aag.json",
    "AI": "spells-ai.json",
    "AitFR-AVT": "spells-aitfr-avt.json",
    "BMT": "spells-bmt.json",
    "DoDk": "spells-dodk.json",
    "EGW": "spells-egw.json",
    "FTD": "spells-ftd.json",
    "GGR": "spells-ggr.json",
    "GHLoE": "spells-ghloe.json",
    "HWCS": "spells-hwcs.json",
    "IDRotF": "spells-idrotf.json",
    "LLK": "spells-llk.json",
    "PHB": "spells-phb.json",
    "SatO": "spells-sato.json",
    "SCC": "spells-scc.json",
    "TCE": "spells-tce.json",
    "TDCSR": "spells-tdcsr.json",
    "XGE": "spells-xge.json",
    "XPHB": "spells-xphb.json"
}
```

### spells-phb.json

Only a single spell is included below to save space. The actual file contains
many more.

```json
{
    "spell": [
        {
            "name": "Acid Splash",
            "source": "PHB",
            "page": 211,
            "srd": true,
            "basicRules": true,
            "reprintedAs": [
                "Acid Splash|XPHB"
            ],
            "level": 0,
            "school": "C",
            "time": [
                {
                    "number": 1,
                    "unit": "action"
                }
            ],
            "range": {
                "type": "point",
                "distance": {
                    "type": "feet",
                    "amount": 60
                }
            },
            "components": {
                "v": true,
                "s": true
            },
            "duration": [
                {
                    "type": "instant"
                }
            ],
            "entries": [
                "You hurl a bubble of acid. Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take {@damage 1d6} acid damage.",
                "This spell's damage increases by {@dice 1d6} when you reach 5th level ({@damage 2d6}), 11th level ({@damage 3d6}), and 17th level ({@damage 4d6})."
            ],
            "scalingLevelDice": {
                "label": "acid damage",
                "scaling": {
                    "1": "1d6",
                    "5": "2d6",
                    "11": "3d6",
                    "17": "4d6"
                }
            },
            "damageInflict": [
                "acid"
            ],
            "savingThrow": [
                "dexterity"
            ],
            "miscTags": [
                "SCL",
                "SGT"
            ],
            "areaTags": [
                "MT",
                "ST"
            ]
        }
    ]
}
```

## Data Directory Layout

Below are the output of the ``ls -R`` command for the data directory. 

```bash
PS C:\src\github\5etools-src\data> ls -R


    Directory: C:\src\github\5etools-src\data


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
d-----         9/28/2024   7:04 PM                adventure
d-----         9/28/2024   7:04 PM                bestiary
d-----         9/28/2024   7:04 PM                book
d-----         9/28/2024   7:04 PM                class
d-----         9/28/2024   7:04 PM                generated
d-----         2/23/2025   5:36 PM                spells
-a----         9/28/2024  11:06 PM          45045 actions.json
-a----         9/28/2024  11:06 PM         701851 adventures.json
-a----         9/28/2024  11:06 PM         865722 backgrounds.json
-a----         9/28/2024  11:06 PM         117270 books.json
-a----         9/28/2024  11:06 PM         272780 changelog.json
-a----         9/28/2024  11:06 PM          69884 charcreationoptions.json
-a----         9/28/2024  11:06 PM          53188 conditionsdiseases.json
-a----         9/28/2024  11:06 PM          73110 cultsboons.json
-a----         9/28/2024  11:06 PM         600823 decks.json
-a----         9/28/2024  11:06 PM         754153 deities.json
-a----         9/28/2024  11:06 PM         305809 encounters.json
-a----         9/28/2024  11:06 PM         224374 feats.json
-a----         9/28/2024  11:06 PM         154831 fluff-backgrounds.json
-a----         9/28/2024  11:06 PM           1888 fluff-charcreationoptions.json
-a----         9/28/2024  11:06 PM           2726 fluff-conditionsdiseases.json
-a----         9/28/2024  11:06 PM           4551 fluff-feats.json
-a----         9/28/2024  11:06 PM         147301 fluff-items.json
-a----         9/28/2024  11:06 PM           5272 fluff-languages.json
-a----         9/28/2024  11:06 PM           2290 fluff-objects.json
-a----         9/28/2024  11:06 PM            258 fluff-optionalfeatures.json
-a----         9/28/2024  11:06 PM         483834 fluff-races.json
-a----         9/28/2024  11:06 PM         246537 fluff-recipes.json
-a----         9/28/2024  11:06 PM            439 fluff-rewards.json
-a----         9/28/2024  11:06 PM            907 fluff-trapshazards.json
-a----         9/28/2024  11:06 PM          61432 fluff-vehicles.json
-a----         9/28/2024  11:06 PM           1287 foundry-actions.json
-a----         9/28/2024  11:06 PM           2611 foundry-feats.json
-a----         9/28/2024  11:06 PM           3627 foundry-items.json
-a----         9/28/2024  11:06 PM          39871 foundry-optionalfeatures.json
-a----         9/28/2024  11:06 PM              3 foundry-psionics.json
-a----         9/28/2024  11:06 PM          24688 foundry-races.json
-a----         9/28/2024  11:06 PM            425 foundry-rewards.json
-a----         9/28/2024  11:06 PM            406 foundry-vehicles.json
-a----         9/28/2024  11:06 PM         109901 items-base.json
-a----         9/28/2024  11:06 PM        2249090 items.json
-a----         9/28/2024  11:06 PM          56491 languages.json
-a----         9/28/2024  11:06 PM          54454 life.json
-a----         9/28/2024  11:06 PM          81915 loot.json
-a----         9/28/2024  11:06 PM         131126 magicvariants.json
-a----         9/28/2024  11:06 PM         111178 makebrew-creature.json
-a----         9/28/2024  11:06 PM           5324 makecards.json
-a----         9/28/2024  11:06 PM           6186 monsterfeatures.json
-a----         9/28/2024  11:06 PM           5298 msbcr.json
-a----         9/28/2024  11:06 PM         230277 names.json
-a----         9/28/2024  11:06 PM          24881 objects.json
-a----         9/28/2024  11:06 PM         145947 optionalfeatures.json
-a----         9/28/2024  11:06 PM         108963 psionics.json
-a----         9/28/2024  11:06 PM         556208 races.json
-a----         9/28/2024  11:06 PM         722595 recipes.json
-a----         9/28/2024  11:06 PM          39112 renderdemo.json
-a----         9/28/2024  11:06 PM         109252 rewards.json
-a----         9/28/2024  11:06 PM           4826 senses.json
-a----         9/28/2024  11:06 PM          12579 skills.json
-a----         9/28/2024  11:06 PM          27307 tables.json
-a----         9/28/2024  11:06 PM          92111 trapshazards.json
-a----         9/28/2024  11:06 PM         699200 variantrules.json
-a----         9/28/2024  11:06 PM          88259 vehicles.json


    Directory: C:\src\github\5etools-src\data\adventure


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         9/28/2024  11:06 PM          40944 adventure-aitfr-avt.json
-a----         9/28/2024  11:06 PM          54163 adventure-aitfr-dn.json
-a----         9/28/2024  11:06 PM          48842 adventure-aitfr-fcd.json
-a----         9/28/2024  11:06 PM          60204 adventure-aitfr-isf.json
-a----         9/28/2024  11:06 PM          57134 adventure-aitfr-thp.json
-a----         9/28/2024  11:06 PM          51056 adventure-azfyt.json
-a----         9/28/2024  11:06 PM        1295482 adventure-bgdia.json
-a----         9/28/2024  11:06 PM        1369085 adventure-cm.json
-a----         9/28/2024  11:06 PM         991854 adventure-coa.json
-a----         9/28/2024  11:06 PM        1827197 adventure-cos.json
-a----         9/28/2024  11:06 PM        1093787 adventure-crcotn.json
-a----         9/28/2024  11:06 PM         301444 adventure-dc.json
-a----         9/28/2024  11:06 PM         188605 adventure-dd.json
-a----         9/28/2024  11:06 PM         711189 adventure-dip.json
-a----         9/28/2024  11:06 PM         178831 adventure-ditlcot.json
-a----         9/28/2024  11:06 PM        1546578 adventure-dodk.json
-a----         9/28/2024  11:06 PM         284324 adventure-dosi.json
-a----         9/28/2024  11:06 PM        1115815 adventure-dsotdq.json
-a----         9/28/2024  11:06 PM         101061 adventure-efr.json
-a----         9/28/2024  11:06 PM         137360 adventure-fs.json
-a----         9/28/2024  11:06 PM         819597 adventure-ghloe.json
-a----         9/28/2024  11:06 PM        1484470 adventure-gos.json
-a----         9/28/2024  11:06 PM          79207 adventure-gotsf.json
-a----         9/28/2024  11:06 PM          61332 adventure-hfstcm.json
-a----         9/28/2024  11:06 PM          86030 adventure-hftt.json
-a----         9/28/2024  11:06 PM         172439 adventure-hol.json
-a----         9/28/2024  11:06 PM         743512 adventure-hotdq.json
-a----         9/28/2024  11:06 PM         444869 adventure-hwaitw.json
-a----         9/28/2024  11:06 PM        2181142 adventure-idrotf.json
-a----         9/28/2024  11:06 PM         377608 adventure-imr.json
-a----         9/28/2024  11:06 PM        1389484 adventure-jttrc.json
-a----         9/28/2024  11:06 PM        1176620 adventure-kftgv.json
-a----         9/28/2024  11:06 PM          61284 adventure-kkw.json
-a----         9/28/2024  11:06 PM          50971 adventure-lk.json
-a----         9/28/2024  11:06 PM         262225 adventure-llk.json
-a----         9/28/2024  11:06 PM         344071 adventure-lmop.json
-a----         9/28/2024  11:06 PM         351434 adventure-lox.json
-a----         9/28/2024  11:06 PM         130557 adventure-lr.json
-a----         9/28/2024  11:06 PM          78243 adventure-lrdt.json
-a----         9/28/2024  11:06 PM          98730 adventure-mot-nss.json
-a----         9/28/2024  11:06 PM          28769 adventure-nrh-ass.json
-a----         9/28/2024  11:06 PM          48818 adventure-nrh-at.json
-a----         9/28/2024  11:06 PM          39160 adventure-nrh-avitw.json
-a----         9/28/2024  11:06 PM          24765 adventure-nrh-awol.json
-a----         9/28/2024  11:06 PM          39936 adventure-nrh-coi.json
-a----         9/28/2024  11:06 PM          49433 adventure-nrh-tcmc.json
-a----         9/28/2024  11:06 PM          34701 adventure-nrh-tlt.json
-a----         9/28/2024  11:06 PM        1752287 adventure-oota.json
-a----         9/28/2024  11:06 PM         789659 adventure-oow.json
-a----         9/28/2024  11:06 PM        1481450 adventure-pabtso.json
-a----         9/28/2024  11:06 PM          79773 adventure-pip.json
-a----         9/28/2024  11:06 PM        1879304 adventure-pota.json
-a----         9/28/2024  11:06 PM        1490144 adventure-qftis.json
-a----         9/28/2024  11:06 PM         241048 adventure-rmbre.json
-a----         9/28/2024  11:06 PM         686901 adventure-rot.json
-a----         9/28/2024  11:06 PM         256659 adventure-rtg.json
-a----         9/28/2024  11:06 PM         204765 adventure-scc-arir.json
-a----         9/28/2024  11:06 PM         247704 adventure-scc-ck.json
-a----         9/28/2024  11:06 PM         156018 adventure-scc-hfmt.json
-a----         9/28/2024  11:06 PM         186742 adventure-scc-tmm.json
-a----         9/28/2024  11:06 PM         369389 adventure-sdw.json
-a----         9/28/2024  11:06 PM         219956 adventure-sja.json
-a----         9/28/2024  11:06 PM        1959399 adventure-skt.json
-a----         9/28/2024  11:06 PM         269542 adventure-slw.json
-a----         9/28/2024  11:06 PM         868422 adventure-tftyp-atg.json
-a----         9/28/2024  11:06 PM         912144 adventure-tftyp-dit.json
-a----         9/28/2024  11:06 PM         384156 adventure-tftyp-tfof.json
-a----         9/28/2024  11:06 PM         289391 adventure-tftyp-thsot.json
-a----         9/28/2024  11:06 PM         176994 adventure-tftyp-toh.json
-a----         9/28/2024  11:06 PM         299737 adventure-tftyp-tsc.json
-a----         9/28/2024  11:06 PM         155593 adventure-tftyp-wpm.json
-a----         9/28/2024  11:06 PM          31100 adventure-tlk.json
-a----         9/28/2024  11:06 PM        1869418 adventure-toa.json
-a----         9/28/2024  11:06 PM         582711 adventure-tofw.json
-a----         9/28/2024  11:06 PM          87014 adventure-tor.json
-a----         9/28/2024  11:06 PM         143261 adventure-ttp.json
-a----         9/28/2024  11:06 PM         105339 adventure-us.json
-a----         9/28/2024  11:06 PM          63971 adventure-uthftlh.json
-a----         9/28/2024  11:06 PM        1338394 adventure-veor.json
-a----         9/28/2024  11:06 PM          77234 adventure-vnotee.json
-a----         9/28/2024  11:06 PM        1250068 adventure-wbtw.json
-a----         9/28/2024  11:06 PM        1443787 adventure-wdh.json
-a----         9/28/2024  11:06 PM        2631895 adventure-wdmm.json
-a----         9/28/2024  11:06 PM          45144 adventure-xmts.json


    Directory: C:\src\github\5etools-src\data\bestiary


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         9/28/2024  11:06 PM          12242 bestiary-aatm.json
-a----         9/28/2024  11:06 PM          68718 bestiary-ai.json
-a----         9/28/2024  11:06 PM          11483 bestiary-aitfr-dn.json
-a----         9/28/2024  11:06 PM           7699 bestiary-aitfr-fcd.json
-a----         9/28/2024  11:06 PM           5830 bestiary-aitfr-isf.json
-a----         9/28/2024  11:06 PM          10947 bestiary-aitfr-thp.json
-a----         9/28/2024  11:06 PM          17626 bestiary-awm.json
-a----         9/28/2024  11:06 PM         188777 bestiary-bam.json
-a----         9/28/2024  11:06 PM          98456 bestiary-bgdia.json
-a----         9/28/2024  11:06 PM         229931 bestiary-bgg.json
-a----         9/28/2024  11:06 PM         142228 bestiary-bmt.json
-a----         9/28/2024  11:06 PM         109126 bestiary-cm.json
-a----         9/28/2024  11:06 PM         221432 bestiary-coa.json
-a----         9/28/2024  11:06 PM         105850 bestiary-cos.json
-a----         9/28/2024  11:06 PM         142263 bestiary-crcotn.json
-a----         9/28/2024  11:06 PM          17269 bestiary-dc.json
-a----         9/28/2024  11:06 PM          11616 bestiary-dip.json
-a----         9/28/2024  11:06 PM           3517 bestiary-ditlcot.json
-a----         9/28/2024  11:06 PM           3306 bestiary-dmg.json
-a----         9/28/2024  11:06 PM           4484 bestiary-dod.json
-a----         9/28/2024  11:06 PM         148916 bestiary-dodk.json
-a----         9/28/2024  11:06 PM          12949 bestiary-dosi.json
-a----         9/28/2024  11:06 PM         214887 bestiary-dsotdq.json
-a----         9/28/2024  11:06 PM          92296 bestiary-egw.json
-a----         9/28/2024  11:06 PM         129840 bestiary-erlw.json
-a----         9/28/2024  11:06 PM           6831 bestiary-esk.json
-a----         9/28/2024  11:06 PM         526249 bestiary-ftd.json
-a----         9/28/2024  11:06 PM         228442 bestiary-ggr.json
-a----         9/28/2024  11:06 PM         273589 bestiary-ghloe.json
-a----         9/28/2024  11:06 PM         122521 bestiary-gos.json
-a----         9/28/2024  11:06 PM            772 bestiary-gotsf.json
-a----         9/28/2024  11:06 PM          24226 bestiary-hat-tg.json
-a----         9/28/2024  11:06 PM           6933 bestiary-hftt.json
-a----         9/28/2024  11:06 PM           6769 bestiary-hol.json
-a----         9/28/2024  11:06 PM          51148 bestiary-hotdq.json
-a----         9/28/2024  11:06 PM         184065 bestiary-hwcs.json
-a----         9/28/2024  11:06 PM         158430 bestiary-idrotf.json
-a----         9/28/2024  11:06 PM          49570 bestiary-imr.json
-a----         9/28/2024  11:06 PM          61056 bestiary-jttrc.json
-a----         9/28/2024  11:06 PM          52346 bestiary-kftgv.json
-a----         9/28/2024  11:06 PM           5658 bestiary-kkw.json
-a----         9/28/2024  11:06 PM          23335 bestiary-llk.json
-a----         9/28/2024  11:06 PM          14035 bestiary-lmop.json
-a----         9/28/2024  11:06 PM           8287 bestiary-lox.json
-a----         9/28/2024  11:06 PM          24303 bestiary-lr.json
-a----         9/28/2024  11:06 PM           4697 bestiary-lrdt.json
-a----         9/28/2024  11:06 PM         223989 bestiary-mabjov.json
-a----         9/28/2024  11:06 PM          31777 bestiary-mcv1sc.json
-a----         9/28/2024  11:06 PM          39655 bestiary-mcv2dc.json
-a----         9/28/2024  11:06 PM          13919 bestiary-mcv3mc.json
-a----         9/28/2024  11:06 PM          77031 bestiary-mcv4ec.json
-a----         9/28/2024  11:06 PM          43037 bestiary-mff.json
-a----         9/28/2024  11:06 PM          13258 bestiary-mgelft.json
-a----         9/28/2024  11:06 PM          21285 bestiary-mismv1.json
-a----         9/28/2024  11:06 PM        1631710 bestiary-mm.json
-a----         9/28/2024  11:06 PM         144554 bestiary-mot.json
-a----         9/28/2024  11:06 PM         894141 bestiary-mpmm.json
-a----         9/28/2024  11:06 PM         174279 bestiary-mpp.json
-a----         9/28/2024  11:06 PM         532706 bestiary-mtf.json
-a----         9/28/2024  11:06 PM           1294 bestiary-nrh-ass.json
-a----         9/28/2024  11:06 PM           8047 bestiary-nrh-at.json
-a----         9/28/2024  11:06 PM           2760 bestiary-nrh-avitw.json
-a----         9/28/2024  11:06 PM            529 bestiary-nrh-awol.json
-a----         9/28/2024  11:06 PM           3785 bestiary-nrh-coi.json
-a----         9/28/2024  11:06 PM           1908 bestiary-nrh-tcmc.json
-a----         9/28/2024  11:06 PM            106 bestiary-nrh-tlt.json
-a----         9/28/2024  11:06 PM         115760 bestiary-oota.json
-a----         9/28/2024  11:06 PM          11678 bestiary-oow.json
-a----         9/28/2024  11:06 PM         101265 bestiary-pabtso.json
-a----         9/28/2024  11:06 PM           7360 bestiary-phb.json
-a----         9/28/2024  11:06 PM         131171 bestiary-pota.json
-a----         9/28/2024  11:06 PM          14654 bestiary-ps-a.json
-a----         9/28/2024  11:06 PM           6399 bestiary-ps-d.json
-a----         9/28/2024  11:06 PM          49272 bestiary-ps-i.json
-a----         9/28/2024  11:06 PM          16030 bestiary-ps-k.json
-a----         9/28/2024  11:06 PM          40348 bestiary-ps-x.json
-a----         9/28/2024  11:06 PM          23213 bestiary-ps-z.json
-a----         9/28/2024  11:06 PM         114445 bestiary-qftis.json
-a----         9/28/2024  11:06 PM           8710 bestiary-rmbre.json
-a----         9/28/2024  11:06 PM          39273 bestiary-rot.json
-a----         9/28/2024  11:06 PM           8476 bestiary-rtg.json
-a----         9/28/2024  11:06 PM          10147 bestiary-sads.json
-a----         9/28/2024  11:06 PM         140119 bestiary-scc.json
-a----         9/28/2024  11:06 PM          20117 bestiary-sdw.json
-a----         9/28/2024  11:06 PM         172375 bestiary-skt.json
-a----         9/28/2024  11:06 PM          15732 bestiary-slw.json
-a----         9/28/2024  11:06 PM          59040 bestiary-tce.json
-a----         9/28/2024  11:06 PM         180291 bestiary-tdcsr.json
-a----         9/28/2024  11:06 PM         147297 bestiary-tftyp.json
-a----         9/28/2024  11:06 PM         146740 bestiary-toa.json
-a----         9/28/2024  11:06 PM        1441217 bestiary-tob1-2023.json
-a----         9/28/2024  11:06 PM          10716 bestiary-tofw.json
-a----         9/28/2024  11:06 PM           9712 bestiary-ttp.json
-a----         9/28/2024  11:06 PM           7279 bestiary-vd.json
-a----         9/28/2024  11:06 PM         143201 bestiary-veor.json
-a----         9/28/2024  11:06 PM         484352 bestiary-vgm.json
-a----         9/28/2024  11:06 PM         122944 bestiary-vrgr.json
-a----         9/28/2024  11:06 PM         127306 bestiary-wbtw.json
-a----         9/28/2024  11:06 PM         191339 bestiary-wdh.json
-a----         9/28/2024  11:06 PM         127517 bestiary-wdmm.json
-a----         9/28/2024  11:06 PM           2950 bestiary-xge.json
-a----         9/28/2024  11:06 PM          11588 bestiary-xmm.json
-a----         9/28/2024  11:06 PM         116638 bestiary-xphb.json
-a----         9/28/2024  11:06 PM           5474 fluff-bestiary-aatm.json
-a----         9/28/2024  11:06 PM          36539 fluff-bestiary-ai.json
-a----         9/28/2024  11:06 PM           6234 fluff-bestiary-aitfr-dn.json
-a----         9/28/2024  11:06 PM           4639 fluff-bestiary-aitfr-fcd.json
-a----         9/28/2024  11:06 PM           1241 fluff-bestiary-aitfr-isf.json
-a----         9/28/2024  11:06 PM           4117 fluff-bestiary-aitfr-thp.json
-a----         9/28/2024  11:06 PM           6978 fluff-bestiary-awm.json
-a----         9/28/2024  11:06 PM          84224 fluff-bestiary-bam.json
-a----         9/28/2024  11:06 PM          42613 fluff-bestiary-bgdia.json
-a----         9/28/2024  11:06 PM          93696 fluff-bestiary-bgg.json
-a----         9/28/2024  11:06 PM          46149 fluff-bestiary-bmt.json
-a----         9/28/2024  11:06 PM          29180 fluff-bestiary-cm.json
-a----         9/28/2024  11:06 PM          77252 fluff-bestiary-coa.json
-a----         9/28/2024  11:06 PM          70466 fluff-bestiary-cos.json
-a----         9/28/2024  11:06 PM          83716 fluff-bestiary-crcotn.json
-a----         9/28/2024  11:06 PM            590 fluff-bestiary-dc.json
-a----         9/28/2024  11:06 PM           3249 fluff-bestiary-dip.json
-a----         9/28/2024  11:06 PM           1264 fluff-bestiary-ditlcot.json
-a----         9/28/2024  11:06 PM            598 fluff-bestiary-dmg.json
-a----         9/28/2024  11:06 PM           1887 fluff-bestiary-dod.json
-a----         9/28/2024  11:06 PM          48602 fluff-bestiary-dodk.json
-a----         9/28/2024  11:06 PM           4084 fluff-bestiary-dosi.json
-a----         9/28/2024  11:06 PM          40276 fluff-bestiary-dsotdq.json
-a----         9/28/2024  11:06 PM          41912 fluff-bestiary-egw.json
-a----         9/28/2024  11:06 PM          79649 fluff-bestiary-erlw.json
-a----         9/28/2024  11:06 PM         297553 fluff-bestiary-ftd.json
-a----         9/28/2024  11:06 PM         131890 fluff-bestiary-ggr.json
-a----         9/28/2024  11:06 PM         192186 fluff-bestiary-ghloe.json
-a----         9/28/2024  11:06 PM          36534 fluff-bestiary-gos.json
-a----         9/28/2024  11:06 PM           6330 fluff-bestiary-hat-tg.json
-a----         9/28/2024  11:06 PM           1636 fluff-bestiary-hftt.json
-a----         9/28/2024  11:06 PM           4360 fluff-bestiary-hotdq.json
-a----         9/28/2024  11:06 PM          72678 fluff-bestiary-hwcs.json
-a----         9/28/2024  11:06 PM          90544 fluff-bestiary-idrotf.json
-a----         9/28/2024  11:06 PM           7043 fluff-bestiary-imr.json
-a----         9/28/2024  11:06 PM          27846 fluff-bestiary-jttrc.json
-a----         9/28/2024  11:06 PM          16599 fluff-bestiary-kftgv.json
-a----         9/28/2024  11:06 PM            422 fluff-bestiary-kkw.json
-a----         9/28/2024  11:06 PM           1506 fluff-bestiary-llk.json
-a----         9/28/2024  11:06 PM           3885 fluff-bestiary-lmop.json
-a----         9/28/2024  11:06 PM           3829 fluff-bestiary-lox.json
-a----         9/28/2024  11:06 PM            505 fluff-bestiary-lr.json
-a----         9/28/2024  11:06 PM            453 fluff-bestiary-lrdt.json
-a----         9/28/2024  11:06 PM         157859 fluff-bestiary-mabjov.json
-a----         9/28/2024  11:06 PM          16878 fluff-bestiary-mcv1sc.json
-a----         9/28/2024  11:06 PM          14491 fluff-bestiary-mcv2dc.json
-a----         9/28/2024  11:06 PM           6941 fluff-bestiary-mcv3mc.json
-a----         9/28/2024  11:06 PM          38774 fluff-bestiary-mcv4ec.json
-a----         9/28/2024  11:06 PM          51933 fluff-bestiary-mff.json
-a----         9/28/2024  11:06 PM           4341 fluff-bestiary-mgelft.json
-a----         9/28/2024  11:06 PM           6044 fluff-bestiary-mismv1.json
-a----         9/28/2024  11:06 PM         771945 fluff-bestiary-mm.json
-a----         9/28/2024  11:06 PM          70306 fluff-bestiary-mot.json
-a----         9/28/2024  11:06 PM         368260 fluff-bestiary-mpmm.json
-a----         9/28/2024  11:06 PM          61612 fluff-bestiary-mpp.json
-a----         9/28/2024  11:06 PM         298108 fluff-bestiary-mtf.json
-a----         9/28/2024  11:06 PM          34732 fluff-bestiary-oota.json
-a----         9/28/2024  11:06 PM           3111 fluff-bestiary-oow.json
-a----         9/28/2024  11:06 PM          31429 fluff-bestiary-pabtso.json
-a----         9/28/2024  11:06 PM          46874 fluff-bestiary-pota.json
-a----         9/28/2024  11:06 PM          22582 fluff-bestiary-ps-a.json
-a----         9/28/2024  11:06 PM           6420 fluff-bestiary-ps-d.json
-a----         9/28/2024  11:06 PM          30403 fluff-bestiary-ps-i.json
-a----         9/28/2024  11:06 PM          24204 fluff-bestiary-ps-k.json
-a----         9/28/2024  11:06 PM          30183 fluff-bestiary-ps-x.json
-a----         9/28/2024  11:06 PM          32753 fluff-bestiary-ps-z.json
-a----         9/28/2024  11:06 PM          42949 fluff-bestiary-qftis.json
-a----         9/28/2024  11:06 PM           5513 fluff-bestiary-rmbre.json
-a----         9/28/2024  11:06 PM           2297 fluff-bestiary-rot.json
-a----         9/28/2024  11:06 PM           7611 fluff-bestiary-sads.json
-a----         9/28/2024  11:06 PM          53954 fluff-bestiary-scc.json
-a----         9/28/2024  11:06 PM           1690 fluff-bestiary-sdw.json
-a----         9/28/2024  11:06 PM         128324 fluff-bestiary-skt.json
-a----         9/28/2024  11:06 PM           1093 fluff-bestiary-tce.json
-a----         9/28/2024  11:06 PM         114389 fluff-bestiary-tdcsr.json
-a----         9/28/2024  11:06 PM          61042 fluff-bestiary-tftyp.json
-a----         9/28/2024  11:06 PM          70074 fluff-bestiary-toa.json
-a----         9/28/2024  11:06 PM         857579 fluff-bestiary-tob1-2023.json
-a----         9/28/2024  11:06 PM           3915 fluff-bestiary-tofw.json
-a----         9/28/2024  11:06 PM           5319 fluff-bestiary-ttp.json
-a----         9/28/2024  11:06 PM           5126 fluff-bestiary-vd.json
-a----         9/28/2024  11:06 PM          46730 fluff-bestiary-veor.json
-a----         9/28/2024  11:06 PM         327741 fluff-bestiary-vgm.json
-a----         9/28/2024  11:06 PM          80963 fluff-bestiary-vrgr.json
-a----         9/28/2024  11:06 PM          41869 fluff-bestiary-wbtw.json
-a----         9/28/2024  11:06 PM          78466 fluff-bestiary-wdh.json
-a----         9/28/2024  11:06 PM          27325 fluff-bestiary-wdmm.json
-a----         9/28/2024  11:06 PM          10594 fluff-bestiary-xphb.json
-a----         9/28/2024  11:06 PM           3287 fluff-index.json
-a----         9/28/2024  11:06 PM           4358 foundry.json
-a----         9/28/2024  11:06 PM           3260 index.json
-a----         9/28/2024  11:06 PM         400264 legendarygroups.json
-a----         9/28/2024  11:06 PM          72254 template.json


    Directory: C:\src\github\5etools-src\data\book


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         9/28/2024  11:06 PM          95918 book-aag.json
-a----         9/28/2024  11:06 PM          98251 book-aatm.json
-a----         9/28/2024  11:06 PM         433599 book-ai.json
-a----         9/28/2024  11:06 PM          27923 book-al.json
-a----         9/28/2024  11:06 PM          41259 book-bam.json
-a----         9/28/2024  11:06 PM         479963 book-bgg.json
-a----         9/28/2024  11:06 PM         876897 book-bmt.json
-a----         9/28/2024  11:06 PM        1279813 book-dmg.json
-a----         9/28/2024  11:06 PM         237947 book-dmtcrg.json
-a----         9/28/2024  11:06 PM          99656 book-dod.json
-a----         9/28/2024  11:06 PM        1088845 book-egw.json
-a----         9/28/2024  11:06 PM        1318429 book-erlw.json
-a----         9/28/2024  11:06 PM         745501 book-ftd.json
-a----         9/28/2024  11:06 PM         698567 book-ggr.json
-a----         9/28/2024  11:06 PM           4272 book-hat-tg.json
-a----         9/28/2024  11:06 PM         188572 book-hf.json
-a----         9/28/2024  11:06 PM         233317 book-hffotm.json
-a----         9/28/2024  11:06 PM         204058 book-hwcs.json
-a----         9/28/2024  11:06 PM         233431 book-mabjov.json
-a----         9/28/2024  11:06 PM          14210 book-mcv4ec.json
-a----         9/28/2024  11:06 PM          86554 book-mm.json
-a----         9/28/2024  11:06 PM         646190 book-mot.json
-a----         9/28/2024  11:06 PM          30757 book-mpmm.json
-a----         9/28/2024  11:06 PM          91169 book-mpp.json
-a----         9/28/2024  11:06 PM         600187 book-mtf.json
-a----         9/28/2024  11:06 PM           8531 book-oga.json
-a----         9/28/2024  11:06 PM         143867 book-paf.json
-a----         9/28/2024  11:06 PM         557747 book-phb.json
-a----         9/28/2024  11:06 PM          33685 book-ps-a.json
-a----         9/28/2024  11:06 PM          94534 book-ps-d.json
-a----         9/28/2024  11:06 PM          74504 book-ps-i.json
-a----         9/28/2024  11:06 PM          39571 book-ps-k.json
-a----         9/28/2024  11:06 PM          89096 book-ps-x.json
-a----         9/28/2024  11:06 PM          17343 book-ps-z.json
-a----         9/28/2024  11:06 PM         283872 book-rmr.json
-a----         9/28/2024  11:06 PM         185905 book-sac.json
-a----         9/28/2024  11:06 PM         421395 book-sato.json
-a----         9/28/2024  11:06 PM         753546 book-scag.json
-a----         9/28/2024  11:06 PM         235985 book-scc.json
-a----         9/28/2024  11:06 PM          13593 book-screen.json
-a----         9/28/2024  11:06 PM          46648 book-screendungeonkit.json
-a----         9/28/2024  11:06 PM          29984 book-screenspelljammer.json
-a----         9/28/2024  11:06 PM          22228 book-screenwildernesskit.json
-a----         9/28/2024  11:06 PM         443622 book-tce.json
-a----         9/28/2024  11:06 PM          65485 book-td.json
-a----         9/28/2024  11:06 PM        1006812 book-tdcsr.json
-a----         9/28/2024  11:06 PM          41016 book-tob1-2023.json
-a----         9/28/2024  11:06 PM         562116 book-vgm.json
-a----         9/28/2024  11:06 PM        1010488 book-vrgr.json
-a----         9/28/2024  11:06 PM            131 book-xdmg.json
-a----         9/28/2024  11:06 PM         331293 book-xge.json
-a----         9/28/2024  11:06 PM            131 book-xmm.json
-a----         9/28/2024  11:06 PM         514840 book-xphb.json


    Directory: C:\src\github\5etools-src\data\class


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         9/28/2024  11:06 PM          69283 class-artificer.json
-a----         9/28/2024  11:06 PM         127077 class-barbarian.json
-a----         9/28/2024  11:06 PM         135572 class-bard.json
-a----         9/28/2024  11:06 PM         277542 class-cleric.json
-a----         9/28/2024  11:06 PM         150143 class-druid.json
-a----         9/28/2024  11:06 PM         154022 class-fighter.json
-a----         9/28/2024  11:06 PM         143714 class-monk.json
-a----         9/28/2024  11:06 PM          52190 class-mystic.json
-a----         9/28/2024  11:06 PM         163420 class-paladin.json
-a----         9/28/2024  11:06 PM         131551 class-ranger.json
-a----         9/28/2024  11:06 PM         127564 class-rogue.json
-a----         9/28/2024  11:06 PM          34187 class-sidekick.json
-a----         9/28/2024  11:06 PM         172775 class-sorcerer.json
-a----         9/28/2024  11:06 PM         141544 class-warlock.json
-a----         9/28/2024  11:06 PM         153385 class-wizard.json
-a----         9/28/2024  11:06 PM          10588 fluff-class-artificer.json
-a----         9/28/2024  11:06 PM          18496 fluff-class-barbarian.json
-a----         9/28/2024  11:06 PM          22554 fluff-class-bard.json
-a----         9/28/2024  11:06 PM          21019 fluff-class-cleric.json
-a----         9/28/2024  11:06 PM          46047 fluff-class-druid.json
-a----         9/28/2024  11:06 PM          18334 fluff-class-fighter.json
-a----         9/28/2024  11:06 PM          18113 fluff-class-monk.json
-a----         9/28/2024  11:06 PM           7654 fluff-class-mystic.json
-a----         9/28/2024  11:06 PM          21320 fluff-class-paladin.json
-a----         9/28/2024  11:06 PM          17553 fluff-class-ranger.json
-a----         9/28/2024  11:06 PM          17779 fluff-class-rogue.json
-a----         9/28/2024  11:06 PM           2100 fluff-class-sidekick.json
-a----         9/28/2024  11:06 PM          20946 fluff-class-sorcerer.json
-a----         9/28/2024  11:06 PM          18362 fluff-class-warlock.json
-a----         9/28/2024  11:06 PM          18976 fluff-class-wizard.json
-a----         9/28/2024  11:06 PM            587 fluff-index.json
-a----         9/28/2024  11:06 PM          95274 foundry.json
-a----         9/28/2024  11:06 PM            497 index.json


    Directory: C:\src\github\5etools-src\data\generated


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         9/28/2024  11:06 PM           3879 bookref-dmscreen-index.json
-a----         9/28/2024  11:06 PM         348293 bookref-dmscreen.json
-a----         9/28/2024  11:06 PM         258625 bookref-quick.json
-a----         9/28/2024  11:06 PM        2658238 gendata-maps.json
-a----         9/28/2024  11:06 PM          23860 gendata-nav-adventure-book-index.json
-a----         9/28/2024  11:06 PM         489230 gendata-spell-source-lookup.json
-a----         9/28/2024  11:06 PM          14244 gendata-subclass-lookup.json
-a----         9/28/2024  11:06 PM        2369457 gendata-tables.json
-a----         9/28/2024   7:04 PM          26751 gendata-tag-redirects.json
-a----         9/28/2024  11:06 PM           6782 gendata-variantrules.json
-a----         9/28/2024   7:04 PM            100 README.md


    Directory: C:\src\github\5etools-src\data\spells


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         9/28/2024  11:06 PM            339 fluff-index.json
-a----         9/28/2024  11:06 PM            252 fluff-spells-aag.json
-a----         9/28/2024  11:06 PM           2380 fluff-spells-dodk.json
-a----         9/28/2024  11:06 PM            486 fluff-spells-egw.json
-a----         9/28/2024  11:06 PM            721 fluff-spells-ftd.json
-a----         9/28/2024  11:06 PM            264 fluff-spells-ggr.json
-a----         9/28/2024  11:06 PM           2115 fluff-spells-hwcs.json
-a----         9/28/2024  11:06 PM           2461 fluff-spells-phb.json
-a----         9/28/2024  11:06 PM           3211 fluff-spells-tce.json
-a----         9/28/2024  11:06 PM           2083 fluff-spells-xge.json
-a----         9/28/2024  11:06 PM          12118 fluff-spells-xphb.json
-a----         9/28/2024  11:06 PM          29245 foundry.json
-a----         9/28/2024  11:06 PM            548 index.json
-a----         9/28/2024  11:06 PM         224753 sources.json
-a----         9/28/2024  11:06 PM           1946 spells-aag.json
-a----         2/23/2025   5:36 PM          11201 spells-ai.json
-a----         9/28/2024  11:06 PM           4259 spells-aitfr-avt.json
-a----         9/28/2024  11:06 PM           4192 spells-bmt.json
-a----         9/28/2024  11:06 PM          21262 spells-dodk.json
-a----         9/28/2024  11:06 PM          22932 spells-egw.json
-a----         9/28/2024  11:06 PM          11872 spells-ftd.json
-a----         9/28/2024  11:06 PM           1461 spells-ggr.json
-a----         9/28/2024  11:06 PM           9212 spells-ghloe.json
-a----         9/28/2024  11:06 PM          15410 spells-hwcs.json
-a----         9/28/2024  11:06 PM           2635 spells-idrotf.json
-a----         9/28/2024  11:06 PM           6457 spells-llk.json
-a----         9/28/2024  11:06 PM         596651 spells-phb.json
-a----         9/28/2024  11:06 PM           2352 spells-sato.json
-a----         9/28/2024  11:06 PM           5626 spells-scc.json
-a----         2/23/2025   6:59 PM          34084 spells-tce.json
-a----         9/28/2024  11:06 PM           2769 spells-tdcsr.json
-a----         9/28/2024  11:06 PM         159907 spells-xge.json
-a----         9/28/2024  11:06 PM         572552 spells-xphb.json
´``
