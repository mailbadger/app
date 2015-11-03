<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Campaign extends Model implements Transformable
{
    use TransformableTrait, SoftDeletes;

    protected $table = 'campaigns';

    protected $fillable = [
        'name',
        'subject',
        'from_name',
        'from_email',
        'status',
        'template_id',
        'recipients',
        'sent_at',
    ];

    protected $dates = ['deleted_at', 'sent_at'];

    public function template()
    {
        return $this->belongsTo('newsletters\Entities\Template');
    }

    public function tags()
    {
        return $this->belongsToMany('newsletters\Entities\Tag', 'campaigns_tags', 'campaign_id',
            'tag_id')->withTimestamps();
    }
}
