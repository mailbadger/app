@extends('master')
@section('scripts')
    <script src="{{asset('js/auth.js')}}"></script>
@endsection
@section('content')
    <div class="ui middle aligned center aligned grid">
        <div class="column five wide">
            <h2 class="ui teal image header">
                <div class="content">
                    Log-in to your account
                </div>
            </h2>
            <form class="ui large form" method="POST" action="{{action('Auth\AuthController@postLogin')}}">
                {!! csrf_field() !!}

                <div class="ui stacked segment">
                    <div class="field">
                        <div class="ui left icon input">
                            <i class="user icon"></i>
                            <input type="email" name="email" placeholder="Email address" value="{{ old('email') }}">
                        </div>
                    </div>
                    <div class="field">
                        <div class="ui left icon input">
                            <i class="lock icon"></i>
                            <input type="password" name="password" id="password" placeholder="Password">
                        </div>
                    </div>
                    <button type="submit" class="ui fluid large teal submit button">Login</button>
                </div>
                <div class="ui error message"></div>
            </form>
            @if (count($errors) > 0)
                <div class="ui error message">
                    <ul class="list">
                        @foreach ($errors->all() as $error)
                            <li>{{ $error }}</li>
                        @endforeach
                    </ul>
                </div>
            @endif
            <div class="ui message">
                New to Newsletter-System? User: admin, Password: admin.
            </div>
        </div>
    </div>
@endsection