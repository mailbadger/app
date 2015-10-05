<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.7.15
 * Time: 18:44.
 */

namespace newsletters\Http\Controllers;

use Illuminate\Http\Request;

class DashboardController extends Controller
{
    public function __construct(Request $request)
    {
        $this->middleware('auth');
        view()->share('activeSidebar', last($request->segments()));
    }

    public function getIndex()
    {
        return view('dashboard.campaigns.list');
    }

    public function getNewCampaign()
    {
        return view('dashboard.campaigns.create_new');
    }

    public function getTemplates()
    {
        return view('dashboard.templates.list');
    }

    public function getNewTemplate()
    {
        return view('dashboard.templates.create_new');
    }

    public function getSubscribers()
    {
        return view('dashboard.subscribers.list');
    }

    public function getNewSubscribers()
    {
        return view('dashboard.subscribers.create_new');
    }

    public function getSettings()
    {
        return view('dashboard.settings');
    }
}
