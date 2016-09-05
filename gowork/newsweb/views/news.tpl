<!DOCTYPE html>
<!-- BEGIN HEAD -->

<head>

	<meta charset="utf-8" />

	<title>{{.NewsTypeName}}</title>

	<meta content="width=device-width, initial-scale=1.0" name="viewport" />

	<meta content="" name="description" />

	<meta content="" name="author" />

	<link href="/static/css/bootstrap.min.css" rel="stylesheet" type="text/css"/>
	<link href="/static/css/default.css" rel="stylesheet" type="text/css" id="style_color"/>

</head>

<!-- END HEAD -->

<!-- BEGIN BODY -->

<body>
	<!-- BEGIN CONTAINER -->   

	<div class="page-container">
		<div class="page-content">
			<div class="container-fluid row-fluid">											
							<div class="span11 row-fluid">
								<h3>{{.NewsTitle}}</h3>							
									<ul class="unstyled inline">
											<li id="news_time"><a href="#">时间：{{.NewsDateTime}}</a></li>
											<li id="news_autor"><a href="#">作者：{{.NewsAutor}}</a></li>
									</ul>						
                                <!--p align="left"><img src="{{.NewsImg}}" alt=""></p-->								  
                                    {{str2html .NewsContent}} 

							</div>

							<!--end span9-->


					

					</div>

			

				<!-- END PAGE CONTENT-->

			</div>

			<!-- END PAGE CONTAINER--> 

		</div>

		<!-- END PAGE -->    

	</div>

</body>

</html>