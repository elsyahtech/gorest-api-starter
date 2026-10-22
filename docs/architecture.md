# Architecture and application lifecycle

This guide explains how a request moves through the starter, what each package owns, and the conventions used when adding application code. The implementation follows a layered flow so HTTP concerns, business rules, and database access have clear boundaries.

## End-to-end lifecycle

The lifecycle has two stages: application startup, followed by request processing. The diagram describes this repository's wiring; details such as middleware behavior and database APIs are defined by GOREST and its dependencies.

```drawio
<mxfile>
  <diagram id="YGAEcujr8at42SBnAT23" name="Page-1">
    <mxGraphModel dx="1243" dy="1504" grid="1" gridSize="10" guides="1" tooltips="0" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="850" pageHeight="1100" math="0" shadow="0">
      <root>
        <mxCell id="0" />
        <mxCell id="1" parent="0" />
        <UserObject label="" mermaidData="{&#xa;  &quot;data&quot;: &quot;flowchart TD\n    A[\&quot;main.go: create GOREST app\&quot;] --&gt; B[\&quot;Register logger, Redis, database, server\&quot;]\n    B --&gt; C[\&quot;Register User, Department, Order routers\&quot;]\n    C --&gt; D{\&quot;Startup services succeed?\&quot;}\n    D -- \&quot;No\&quot; --&gt; E[\&quot;Startup fails; inspect configuration and logs\&quot;]\n    D -- \&quot;Yes\&quot; --&gt; F[\&quot;Serve HTTP requests\&quot;]\n    F --&gt; Q{\&quot;SIGINT or SIGTERM received?\&quot;}\n    Q -- \&quot;No\&quot; --&gt; F\n    Q -- \&quot;Yes\&quot; --&gt; R[\&quot;app.Close(): shut down application services\&quot;]\n    F --&gt; G[\&quot;Server middleware&lt;br/&gt;recover, request ID, compression,&lt;br/&gt;CORS, helmet, rate limiter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-50" edge="1" parent="42" source="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;entryX=0.5;entryY=0;entryDx=0;entryDy=0;exitX=0.178;exitY=1.014;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="460.0000000000002" y="670" as="targetPoint" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="386" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="292" y="870" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-51" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.0364" y="1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-53" connectable="0" parent="27" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2384" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```miter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-50" edge="1" parent="42" source="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;entryX=0.5;entryY=0;entryDx=0;entryDy=0;exitX=0.178;exitY=1.014;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="460.0000000000002" y="670" as="targetPoint" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="386" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="292" y="870" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-51" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.0364" y="1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-53" connectable="0" parent="27" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2384" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```iter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-50" edge="1" parent="42" source="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;entryX=0.5;entryY=0;entryDx=0;entryDy=0;exitX=0.178;exitY=1.014;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="460.0000000000002" y="670" as="targetPoint" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="386" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="292" y="870" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-51" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.0364" y="1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Yes" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```iter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-50" edge="1" parent="42" source="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;entryX=0.5;entryY=0;entryDx=0;entryDy=0;exitX=0.178;exitY=1.014;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="460.0000000000002" y="670" as="targetPoint" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="386" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="292" y="870" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-51" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.0364" y="1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Yes" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```iter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-50" edge="1" parent="42" source="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;entryX=0.5;entryY=0;entryDx=0;entryDy=0;exitX=0.178;exitY=1.014;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="460.0000000000002" y="670" as="targetPoint" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="386" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="292" y="870" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-51" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.0364" y="1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Yes" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```iter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-50" edge="1" parent="42" source="7" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;entryX=0.5;entryY=0;entryDx=0;entryDy=0;exitX=0.178;exitY=1.014;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="460.0000000000002" y="670" as="targetPoint" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="386" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="292" y="870" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-46" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2083" y="3" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Yes" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```miter\&quot;]\n    G --&gt; H[\&quot;Router matches method and path\&quot;]\n    H --&gt; I[\&quot;Handler parses path, query, and JSON\&quot;]\n    I --&gt; J[\&quot;Library validates and applies app rules\&quot;]\n    J --&gt; K[\&quot;Model builds ORM or native query\&quot;]\n    K --&gt; L[\&quot;GOREST database API / configured driver\&quot;]\n    L --&gt; M{\&quot;Database operation succeeds?\&quot;}\n    M -- \&quot;Yes\&quot; --&gt; N[\&quot;Handler maps data with app/view&lt;br/&gt;and builds JSON response\&quot;]\n    N --&gt; O[\&quot;HTTP response\&quot;]\n    O --&gt; F\n    I -. \&quot;input or validation error\&quot; .-&gt; P[\&quot;Handler logs error and sends error response\&quot;]\n    J -. \&quot;application or model error\&quot; .-&gt; P\n    M -- \&quot;No\&quot; --&gt; P\n    P --&gt; O&quot;,&#xa;  &quot;config&quot;: null&#xa;}" id="42">
          <mxCell connectable="0" parent="1" style="group;transparentBounds=1;editIcon=1;lockedGroup=0;groupPadding=10;" vertex="1">
            <mxGeometry as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="main.go: create GOREST app" mermaidId="n:A" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="main.go: create GOREST app" id="2">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="267" x="156" y="10" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register logger, Redis, database, server" mermaidId="n:B" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register logger, Redis, database, server" id="3">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="343" x="118" y="114" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Register User, Department, Order routers" mermaidId="n:C" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Register User, Department, Order routers" id="4">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="356" x="112" y="218" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup services succeed?" mermaidId="n:D" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup services succeed?" id="5">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="158" width="158" x="210.5" y="322" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Startup fails; inspect configuration and logs" mermaidId="n:E" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Startup fails; inspect configuration and logs" id="6">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="220" x="-21" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Serve HTTP requests" mermaidId="n:F" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Serve HTTP requests" id="7">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="331" x="400" y="521" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="SIGINT or SIGTERM received?" mermaidId="n:Q" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="SIGINT or SIGTERM received?" id="8">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="147" width="147" x="290" y="680" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="app.Close(): shut down application services" mermaidId="n:R" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="app.Close(): shut down application services" id="9">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="46" width="337" x="193.5" y="909" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" mermaidId="n:G" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Server middleware&#xa;recover, request ID, compression,&#xa;CORS, helmet, rate limiter" id="10">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="93" width="301" x="510" y="1007" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Router matches method and path" mermaidId="n:H" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Router matches method and path" id="11">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1160" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler parses path, query, and JSON" mermaidId="n:I" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler parses path, query, and JSON" id="12">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="299" x="510" y="1270" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Library validates and applies app rules" mermaidId="n:J" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Library validates and applies app rules" id="13">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="404" x="330" y="1430" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Model builds ORM or native query" mermaidId="n:K" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Model builds ORM or native query" id="14">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="300" x="248" y="1548" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="GOREST database API / configured driver" mermaidId="n:L" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="GOREST database API / configured driver" id="15">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="354" x="221" y="1671" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Database operation succeeds?" mermaidId="n:M" mermaidBaseStyle="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Database operation succeeds?" id="16">
          <mxCell parent="42" style="rhombus;html=1;strokeWidth=1;whiteSpace=wrap;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="185" width="185" x="305.5" y="1780" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-48" edge="1" parent="42" source="17" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=0;entryY=0.5;entryDx=0;entryDy=0;" target="18">
          <mxGeometry relative="1" as="geometry" />
        </mxCell>
        <UserObject label="Handler maps data with app/view&#xa;and builds JSON response" mermaidId="n:N" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler maps data with app/view&#xa;and builds JSON response" id="17">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="73" width="304" x="246" y="2050" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-49" edge="1" parent="42" source="18" style="edgeStyle=orthogonalEdgeStyle;rounded=0;orthogonalLoop=1;jettySize=auto;html=1;exitX=1;exitY=0.5;exitDx=0;exitDy=0;entryX=1;entryY=0.5;entryDx=0;entryDy=0;" target="7">
          <mxGeometry relative="1" as="geometry">
            <Array as="points">
              <mxPoint x="1000" y="2086.5" />
              <mxPoint x="1000" y="548" />
            </Array>
          </mxGeometry>
        </mxCell>
        <UserObject label="HTTP response" mermaidId="n:O" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="HTTP response" id="18">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="165" x="688" y="2059.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="Handler logs error and sends error response" mermaidId="n:P" mermaidBaseStyle="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" mermaidBaseValue="Handler logs error and sends error response" id="19">
          <mxCell parent="42" style="html=1;whiteSpace=wrap;strokeWidth=1;fillColor=light-dark(#ECECFF,#1f2020);strokeColor=light-dark(#9370DB,#cccccc);fontColor=light-dark(#333333,#cccccc);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontSize=16;" vertex="1">
            <mxGeometry height="54" width="319" x="611" y="1845.5" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:A-&gt;B#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="20">
          <mxCell edge="1" parent="42" source="2" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="3">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:B-&gt;C#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="21">
          <mxCell edge="1" parent="42" source="3" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="4">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:C-&gt;D#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="22">
          <mxCell edge="1" parent="42" source="4" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="5">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:D-&gt;E#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.26;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="No" id="23">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0;exitY=0.5;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="6">
            <mxGeometry relative="1" x="0.0012" as="geometry">
              <Array as="points">
                <mxPoint x="90" y="401" />
              </Array>
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-43" connectable="0" parent="23" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5538" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:D-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.88;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="24">
          <mxCell edge="1" parent="42" source="5" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="7">
            <mxGeometry relative="1" x="-0.0001" as="geometry" />
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-44" connectable="0" parent="24" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;Yes&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="0.5872" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="" mermaidId="e:F-&gt;Q#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5289855072463768;exitY=1;entryX=0.5174418604651163;entryY=0;curved=1;exitX=0.53;exitY=1;entryX=0.52;entryY=0;" mermaidBaseValue="" id="25">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;exitX=0.271;exitY=1.004;exitDx=0;exitDy=0;exitPerimeter=0;" target="8">
            <mxGeometry relative="1" x="-0.3706" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="489.7" y="610" />
                <mxPoint x="490" y="627.5" />
                <mxPoint x="363.5" y="627.5" />
              </Array>
              <mxPoint x="510" y="610" as="sourcePoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:Q-&gt;F#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.4622093023255814;exitY=0;entryX=0.43719806763285024;entryY=1;curved=1;exitX=0.46;exitY=0;entryX=0.44;entryY=1;" mermaidBaseValue="No" id="26">
          <mxCell edge="1" parent="42" source="8" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.01;entryY=0.567;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="7">
            <mxGeometry relative="1" x="0.2262" y="7.36" as="geometry">
              <mxPoint as="offset" />
              <Array as="points">
                <mxPoint x="250" y="755.5" />
                <mxPoint x="250" y="551.6" />
              </Array>
              <mxPoint x="270" y="755.5" as="sourcePoint" />
              <mxPoint x="380" y="550" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-46" connectable="0" parent="26" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;font style=&quot;font-size: 16px;&quot;&gt;No&lt;/font&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2083" y="3" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="Yes" mermaidId="e:Q-&gt;R#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="27">
          <mxCell edge="1" parent="42" source="8" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="9">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:F-&gt;G#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.43;entryY=0;" mermaidBaseValue="" id="28">
          <mxCell edge="1" parent="42" source="7" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.0069" as="geometry">
              <Array as="points">
                <mxPoint x="654.9" y="1007" />
              </Array>
              <mxPoint x="654.9" y="1010" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:G-&gt;H#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="29">
          <mxCell edge="1" parent="42" source="10" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="11">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:H-&gt;I#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="30">
          <mxCell edge="1" parent="42" source="11" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="12">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:I-&gt;J#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="31">
          <mxCell edge="1" parent="42" source="12" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.19;exitY=1;entryX=0.585;entryY=0.026;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;" target="13">
            <mxGeometry relative="1" x="0.301" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;K#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.25;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="32">
          <mxCell edge="1" parent="42" source="13" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.169;exitY=1.013;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;exitPerimeter=0;" target="14">
            <mxGeometry relative="1" x="-0.3928" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:K-&gt;L#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="33">
          <mxCell edge="1" parent="42" source="14" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="15">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:L-&gt;M#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="" id="34">
          <mxCell edge="1" parent="42" source="15" style="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;" target="16">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="Yes" mermaidId="e:M-&gt;N#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.45;exitY=1;entryX=0.5;entryY=0;" mermaidBaseValue="Yes" id="35">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;exitDx=0;exitDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;" target="17">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="input or validation error" mermaidId="e:I-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.77;exitY=1;entryX=0.58;entryY=0;" mermaidBaseValue="input or validation error" id="38">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);entryX=0.463;entryY=0.009;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;entryPerimeter=0;exitX=0.964;exitY=1.002;exitDx=0;exitDy=0;exitPerimeter=0;">
            <mxGeometry relative="1" x="-0.0425" as="geometry">
              <mxPoint x="775.2359999999999" y="1323.1080000000002" as="sourcePoint" />
              <mxPoint x="774.902" y="1844.9859999999999" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:J-&gt;P#0" mermaidBaseStyle="curved=1;dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;entryX=0.42;entryY=0;" mermaidBaseValue="application or model error" id="39">
          <mxCell edge="1" parent="42" style="dashed=1;dashPattern=2 3;fixDash=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.82;exitY=1;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;">
            <mxGeometry relative="1" x="0.323" as="geometry">
              <Array as="points" />
              <mxPoint x="661.2800000000009" y="1487" as="sourcePoint" />
              <mxPoint x="661.62" y="1844.99" as="targetPoint" />
            </mxGeometry>
          </mxCell>
        </UserObject>
        <mxCell id="tPXL4iwxObKe7BxADoTg-47" connectable="0" parent="39" style="edgeLabel;html=1;align=center;verticalAlign=middle;resizable=0;points=[];" value="&lt;span style=&quot;font-size: 16px;&quot;&gt;application or model error&lt;/span&gt;" vertex="1">
          <mxGeometry relative="1" x="-0.2576" y="-1" as="geometry">
            <mxPoint as="offset" />
          </mxGeometry>
        </mxCell>
        <UserObject label="No" mermaidId="e:M-&gt;P#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=0.96;exitY=1;entryX=0.24;entryY=0;" mermaidBaseValue="No" id="40">
          <mxCell edge="1" parent="42" source="16" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);html=1;fontSize=16;labelBackgroundColor=light-dark(#E8E8E88D,#2a2a2a8D);fontFamily=Trebuchet MS,Verdana,Arial,sans-serif;fontColor=light-dark(#333333,#cccccc);exitX=1;exitY=0.5;entryX=0;entryY=0.5;entryDx=0;entryDy=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;exitDx=0;exitDy=0;" target="19">
            <mxGeometry relative="1" as="geometry" />
          </mxCell>
        </UserObject>
        <UserObject label="" mermaidId="e:P-&gt;O#0" mermaidBaseStyle="curved=1;startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.76;entryY=0;" mermaidBaseValue="" id="41">
          <mxCell edge="1" parent="42" source="19" style="startArrow=none;endArrow=block;endSize=7;strokeColor=light-dark(#333333,#cccccc);exitX=0.5;exitY=1;entryX=0.5;entryY=0;edgeStyle=orthogonalEdgeStyle;rounded=0;curved=0;entryDx=0;entryDy=0;" target="18">
            <mxGeometry relative="1" as="geometry">
              <Array as="points" />
            </mxGeometry>
          </mxCell>
        </UserObject>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>

```

For a successful `POST /examples/departments`, the request path is:

1. `app/router/department_router.go` matches the method and path and calls `CreateDepartment`.
2. The handler in `app/handler/` decodes the JSON body into `entity.Department`.
3. `DepartmentLibrary` in `app/library/` coordinates the operation and applies application-level checks.
4. `DepartmentModel` in `app/model/` uses GOREST's database API to persist the record.
5. The handler builds the response with the GOREST JSON response helper. If the operation fails, the handler logs the error and returns an error response instead.

User and order endpoints follow the same layer boundaries. User examples also demonstrate password processing, bulk and upsert writes, ORM queries, and native database access.

## Startup behavior

`main.go` creates an application named `gorest-api-starter` and registers the logger, Redis, database, and HTTP server using configuration from `app/config`. It then registers `UserRouter`, `DepartmentRouter`, and `OrderRouter`.

- PostgreSQL database integration and the HTTP server are enabled in the checked-in configuration.
- Redis is registered as an optional service but disabled.
- Database migrations and seeders are enabled. They run as part of database startup before the application serves requests; a startup failure must be investigated rather than treated as a successful launch.
- Configuration is defined in Go source. The application does not load settings from `.env` or environment variables automatically.
- The routes are grouped below `/examples`.

See [Getting started](./getting-started.md) for local setup and [Configuration](./configuration.md) before using a different environment.

## Package responsibilities

This starter is more than a CRUD application that happens to run: it is a practical map for organizing an application as it grows. A small CRUD example could be written entirely in `main.go`, but that structure becomes difficult to navigate and maintain when more resources, validation, business rules, integrations, and database operations are added. When all of those concerns accumulate together, even a working application can become hard to change safely.

The layers below give each kind of work a clear home. This is an intentional starting standard for client applications—not structure for structure's sake, and not a claim that every small feature needs every layer. Follow the boundaries as the application grows so new behavior has a predictable place and the codebase remains understandable, testable, and maintainable over the long term.

| Directory | Responsibility |
|---|---|
| `app/config` | Holds the application's infrastructure settings, such as logging, Redis, database connections, and the HTTP server. |
| `app/router` | Defines which HTTP method and URL reach which handler. It also assembles the model, library, and handler for each resource, connecting the layers without putting request processing or business rules in the route definitions. |
| `app/handler` | The HTTP-facing layer, or request boundary. It receives client requests, decodes and checks incoming data, and passes the relevant input to a library. It then turns the result into an HTTP response. Keep transport-specific work here; do not let handlers become a place where business rules and database operations accumulate. |
| `app/library` | The application/service layer for business logic. It applies rules, calculations, and workflows, and coordinates one or more model operations. For example, user write methods process passwords here. This is where behavior belongs when it is part of the business rather than HTTP or database mechanics. |
| `app/model` | The data-access layer between the application and the database. It performs persistence, retrieval, and database-specific ORM or native queries through GOREST. It should not decide HTTP response formats or own business workflows. |
| `app/entity` | Defines the application's data structures (Go `struct` types), including entities and the request, patch, and response shapes used by the examples. These types give data a clear structure as it moves between layers; they do not perform the work of those layers. |
| `app/view` | Transforms model entities into response DTOs intended for clients, keeping the public response shape separate from the database representation. |
| `app/database/migrations` | Contains driver-specific files for creating or updating database schemas or collections. |
| `app/database/seeders` | Contains driver-specific files for inserting sample or initial data. |

A typical request flows from the router to a handler, then to a library, and from there to a model when database access is needed. The result returns through the layers and is formatted by the handler for the client; a view may shape entity data into a response DTO along the way.

By keeping HTTP handling, business behavior, and database access from accumulating in the same functions, these boundaries make changes easier to locate, test, and review. Use the structure as a guide rather than adding empty layers: a feature may not need every layer, but each responsibility it does have should remain in its appropriate place.

## Naming and file organization

Use the same resource and operation names across layers, with a suffix that identifies the layer:

```text
app/handler/user_create_handler.go   -> UserHandler.CreateUser
app/library/user_create_library.go   -> UserLibrary.CreateUser
app/model/user_create_model.go       -> UserModel.CreateUser
```

Use lowercase filenames with underscores and the pattern `<resource>_<operation>_<layer>.go`. Keep shared types and constructors in a package-named file such as `handler.go`, `library.go`, or `model.go`; do not call such a file `interface.go` unless it actually contains interface declarations.

For exported methods, use **verb + resource + specific behavior**, for example:

- `CreateUser`, `GetUserByID`, `UpdateDepartment`, `DeleteOrder`
- `CreateUserWithTransaction`, `GetUsersWithNativeMongoQuery`
- `CreateUsersInBulk`, `PatchUserActiveStatus`

Prefer `Get` for read methods, and name special behavior explicitly. Keep names aligned through handler, library, model, and router. Add a Go doc comment immediately above each exported declaration; describe the actual behavior (including relevant query constraints), not merely repeat the name.

## Entities, request types, and response DTOs

Entity tags define JSON property names and GOREST database mappings. The user and order entities also reference related resources:

- `User.DepartmentID` identifies a department.
- `Order.UserID` identifies a user.
- `UserResponse` omits the password field; use response DTOs and `app/view` mappings rather than returning a password-bearing entity.
- Patch request structs use pointers for editable fields. This distinguishes an omitted value from explicit zero values such as `false`, `0`, or `""`.

When changing an entity, check all dependent layers: request validation, model queries, response mapping, API examples, migration and seeder files for the supported drivers, and tests.

## Driver boundaries

The starter includes ORM and native query examples for selected drivers, but does not guarantee identical behavior across them. The native SQL user examples use PostgreSQL syntax. MongoDB and ScyllaDB use their native APIs. GOREST `BeginTx`, `Commit`, and `Rollback` calls do not provide SQL-style transaction atomicity for MongoDB and ScyllaDB in these examples. See the [database guide](./database.md) before adapting these operations.
