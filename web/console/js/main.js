    
    
    var sys = arbor.ParticleSystem(1000, 1000,1);
    sys.parameters({gravity:true});
    sys.renderer = Renderer("#viewport", sys) ;
    var chainNumber = 1;

    var addComponent = function(v, sys) {

        var root = sys.addNode(v.name,{'color':'red','shape':'dot','label':v.name});

        if (v.paramsInput != undefined){
            $.each(v.paramsInput, function(k1, v1){
                var inParam = sys.addNode(v.name + ".in." + k1, {'color':'green','shape':'dot','label':k1});
                sys.addEdge(root, inParam);
            });
        }

        if (v.paramsOutput != undefined){
            $.each(v.paramsOutput, function(k2, v2){
                var outParam = sys.addNode(v.name + ".out." + k2, {'color':'blue','shape':'dot','label':k2});
                sys.addEdge(root, outParam);
            });         
        }
    }

    function recurAddition(item, parent, weight) {
        console.log(weight);
        if ($.isEmptyObject(item)){
            return;
        } 

        if (typeof item == "string"){
            return;
        }
            
        $.each(item, function(key, value){
            if (!isNaN(key)){
                    key = parent.data.label + key;
            }

            var node = sys.addNode(key + Math.random(), {'color':'red', 'shape':'dot', 'label':key, 'weight': weight});
            recurAddition(value, node, weight + 1);

            if (parent != undefined)
                sys.addEdge(parent, node, {color:'green', directed:true})
        });
    }

    // $.each(compArray, function(k, v) {
    //     addComponent(v, sys);
    // });
    var componentMap = {};
    var canvasShow = false;

    $.ajax({url: "http://localhost:8020/components", success: function(data){
            var components = JSON.parse(data);
            $.each(components, function(k, v){
                componentMap[v.name] = v
                $('.dropdown-menu').append('<li><a>' + v.name  +  '</a></li>');
            })

            //adding handlers to controls
            $('.dropdown-menu li a').click(function(event){
                var compName = event.currentTarget.innerText
                addComponent(componentMap[compName], sys)
            })

            $('.add-chain').click(function(event){
                var chainInfo = sys.renderer.getAddedEdges();
                // var dataInfo = JSON.stringify({ChainId:chainNumber, ChainInfo:chainInfo});
                // if (dataInfo == undefined) return;
                var mapComp = {};
                var edgeInfo;
                var compName;
                $.each(chainInfo, function(index, value){
                    if (compName !== value.src.name.split('.')[0]){
                        edgeInfo = []
                    }
                    compName = value.src.name.split('.')[0];
                    edgeInfo.push({src:value.src.name.split('.')[2], dest:value.dest.name});
                    mapComp[compName] = edgeInfo;
                });
                mapComp.ChainId = chainNumber
                console.log(mapComp)
                var mapCompString = JSON.stringify(mapComp);
                $.post({url: "http://localhost:8020/chain",
                 data:mapCompString, success : function(){
                    var htmlText = $('.chainId').html();
                    $('.chainId').html(htmlText + chainNumber)
                    // chainNumber += 1;
                }});
            })

            $('.new-chain').click(function(event){
            })

            $('.hide-chain').click(function(event){
                if (canvasShow == false){
                    $('#viewport').css('display', 'none');
                    $('.hide-chain').html('Show Chain');
                    canvasShow = true;
                }else {
                    $('#viewport').css('display', 'inline-block');
                    $('.hide-chain').html('Hide Chain');
                    canvasShow = false;
                }
            })
            
        }, error: function(){

        }
    });


