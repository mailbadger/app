<nav class="navbar navbar-inverse navbar-fixed-top">
    <div class="container">
        <div class="navbar-header">
            <button aria-controls="navbar" aria-expanded="false" data-target="#navbar" data-toggle="collapse" class="navbar-toggle collapsed" type="button">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a href="{{url()}}" class="navbar-brand">Newsletters</a>
        </div>
        <div class="navbar-collapse collapse" id="navbar">
            <ul class="nav navbar-nav navbar-right">
                @if(Auth::check())
                    <li><a href="#">Dashboard</a></li>
                    <li><a href="#">Settings</a></li>
                    <li><a href="#">Profile</a></li>
                @else
                    <li><a href="{{url('auth/login')}}">Login</a></li>
                @endif
            </ul>
        </div><!--/.nav-collapse -->
    </div>
</nav>